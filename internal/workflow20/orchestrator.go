package workflow20

import (
	"fmt"
	"time"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/repository"
)

type Orchestrator struct {
	repo     *repository.Repository
	clock    func() time.Time
	sequence int
}

type WorkflowResult struct {
	Entry      domain.JournalEntry
	Resources  []domain.Resource
	Event      domain.DomainEvent
	CloseOrder []string
}

func New(repo *repository.Repository) *Orchestrator {
	return &Orchestrator{repo: repo, clock: func() time.Time { return time.Unix(0, 0).UTC() }}
}

func (o *Orchestrator) WithClock(clock func() time.Time) *Orchestrator {
	if clock != nil {
		o.clock = clock
	}
	return o
}

func (o *Orchestrator) nextID(prefix string) string {
	o.sequence++
	return fmt.Sprintf("%s-%03d", prefix, o.sequence)
}

func (o *Orchestrator) CaptureDraft(entry domain.JournalEntry, resources []domain.Resource) (WorkflowResult, error) {
	if entry.Status == "" {
		entry.Status = domain.StatusDraft
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = o.clock()
	}
	entry.UpdatedAt = o.clock()
	entry.Version++
	entry.Achievements = domain.NormalizeList(entry.Achievements)
	entry.Blockers = domain.NormalizeList(entry.Blockers)
	entry.NextSteps = domain.NormalizeList(entry.NextSteps)
	entry.Tags = domain.NormalizeList(entry.Tags)
	if err := entry.Validate(); err != nil {
		return WorkflowResult{}, err
	}
	if issues := domain.ValidateResourceTree(resources); domain.HasErrors(issues) {
		return WorkflowResult{}, fmt.Errorf("resource validation failed: %v", issues)
	}
	if err := o.repo.SaveEntry(entry); err != nil {
		return WorkflowResult{}, err
	}
	for _, resource := range resources {
		resource.EntryID = entry.ID
		if err := o.repo.SaveResource(resource); err != nil {
			return WorkflowResult{}, err
		}
	}
	event := domain.NewEvent(o.nextID("event"), domain.EventCreated, entry.ID, entry.Author, map[string]string{"status": string(entry.Status)}, o.clock())
	if err := o.repo.SaveEvent(event); err != nil {
		return WorkflowResult{}, err
	}
	return WorkflowResult{Entry: entry, Resources: resources, Event: event}, nil
}

func (o *Orchestrator) Submit(entryID, actor string) (domain.JournalEntry, error) {
	entry, err := o.repo.Entry(entryID)
	if err != nil {
		return entry, err
	}
	if !entry.IsEditable() {
		return entry, fmt.Errorf("%w: only draft entries can be submitted", domain.ErrConflict)
	}
	entry.Status = domain.StatusSubmitted
	entry.UpdatedAt = o.clock()
	entry.Version++
	if err := o.repo.SaveEntry(entry); err != nil {
		return entry, err
	}
	event := domain.NewEvent(o.nextID("event"), domain.EventSubmitted, entryID, actor, nil, o.clock())
	return entry, o.repo.SaveEvent(event)
}

func (o *Orchestrator) SaveResourceStep(entryID string, resources []domain.Resource, recorder CloseRecorder) (WorkflowResult, error) {
	entry, err := o.repo.Entry(entryID)
	if err != nil {
		return WorkflowResult{}, err
	}
	parent, child, err := ensureResourceTree(resources)
	if err != nil {
		return WorkflowResult{}, err
	}
	parent.EntryID = entryID
	child.EntryID = entryID
	parentCloser := NewResourceCloser(&parent, nil, recorder)
	childCloser := NewResourceCloser(&child, parentCloser, recorder)
	session := NewResourceSession(o.repo, entry, parentCloser, childCloser)
	return session.Commit()
}
