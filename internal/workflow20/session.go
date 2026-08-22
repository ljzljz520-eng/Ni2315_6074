package workflow20

import (
	"fmt"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/repository"
)

type ResourceSession struct {
	repo   *repository.Repository
	entry  domain.JournalEntry
	parent *ResourceCloser
	child  *ResourceCloser
}

func NewResourceSession(repo *repository.Repository, entry domain.JournalEntry, parent, child *ResourceCloser) *ResourceSession {
	return &ResourceSession{repo: repo, entry: entry, parent: parent, child: child}
}

func (s *ResourceSession) Commit() (result WorkflowResult, err error) {
	defer func() {
		if flushErr := s.flush(); err == nil && flushErr != nil {
			err = flushErr
		}
	}()
	// Defers run LIFO, so the parent must be registered after the child to
	// close first. Closing the parent before the child is the required
	// dependency order: ResourceCloser.Close clears a child's payload when its
	// parent is still open, so closing the child first would discard the
	// payload before flush persists it.
	defer s.child.Close()
	defer s.parent.Close()
	if s.entry.Status != domain.StatusSubmitted && s.entry.Status != domain.StatusDraft {
		return result, fmt.Errorf("%w: resource step requires draft or submitted entry", domain.ErrConflict)
	}
	result.Entry = s.entry
	return result, nil
}

func (s *ResourceSession) flush() error {
	parent := s.parent.Resource()
	child := s.child.Resource()
	if err := s.repo.SaveResource(parent); err != nil {
		return err
	}
	if err := s.repo.SaveResource(child); err != nil {
		return err
	}
	resultEvent := domain.NewEvent("resource-"+s.entry.ID, domain.EventReviewed, s.entry.ID, s.entry.Author, map[string]string{"parent": parent.ID, "child": child.ID}, s.entry.UpdatedAt)
	return s.repo.SaveEvent(resultEvent)
}
