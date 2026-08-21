package repository

import (
	"sort"
	"strings"

	"independentweeklylog/internal/domain"
)

func (r *Repository) ReportFor(author string) (domain.DeveloperReport, error) {
	entries, err := r.Entries()
	if err != nil {
		return domain.DeveloperReport{}, err
	}
	return domain.BuildDeveloperReport(entries, author), nil
}

func (r *Repository) StatusBreakdown() (domain.StatusBreakdown, error) {
	entries, err := r.Entries()
	if err != nil {
		return domain.StatusBreakdown{}, err
	}
	return domain.CountStatuses(entries), nil
}

func (r *Repository) Tags() ([]string, error) {
	entries, err := r.Entries()
	if err != nil {
		return nil, err
	}
	return domain.MergeTags(entries), nil
}

func (r *Repository) Authors() ([]string, error) {
	entries, err := r.Entries()
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool)
	for _, entry := range entries {
		clean := strings.TrimSpace(entry.Author)
		if clean != "" {
			seen[clean] = true
		}
	}
	authors := make([]string, 0, len(seen))
	for author := range seen {
		authors = append(authors, author)
	}
	sort.Strings(authors)
	return authors, nil
}

func (r *Repository) ReplaceEntry(entry domain.JournalEntry) error {
	current, err := r.Entry(entry.ID)
	if err != nil {
		return err
	}
	if entry.Version <= current.Version {
		return domain.ErrConflict
	}
	return r.SaveEntry(entry)
}

func (r *Repository) AddTag(entryID, tag string) error {
	entry, err := r.Entry(entryID)
	if err != nil {
		return err
	}
	if !entry.IsEditable() {
		return domain.ErrConflict
	}
	entry.Tags = domain.NormalizeList(append(entry.Tags, tag))
	entry.Version++
	return r.SaveEntry(entry)
}

func (r *Repository) AddAchievement(entryID, achievement string) error {
	entry, err := r.Entry(entryID)
	if err != nil {
		return err
	}
	if !entry.IsEditable() {
		return domain.ErrConflict
	}
	entry.Achievements = domain.NormalizeList(append(entry.Achievements, achievement))
	entry.Version++
	return r.SaveEntry(entry)
}

func (r *Repository) AddNextStep(entryID, next string) error {
	entry, err := r.Entry(entryID)
	if err != nil {
		return err
	}
	if !entry.IsEditable() {
		return domain.ErrConflict
	}
	entry.NextSteps = domain.NormalizeList(append(entry.NextSteps, next))
	entry.Version++
	return r.SaveEntry(entry)
}

func (r *Repository) CountByWeek() (map[int]int, error) {
	entries, err := r.Entries()
	if err != nil {
		return nil, err
	}
	counts := make(map[int]int)
	for _, entry := range entries {
		counts[entry.Week]++
	}
	return counts, nil
}

func (r *Repository) Latest(author string) (domain.JournalEntry, error) {
	entries, err := r.FindEntries(EntryFilter{Author: author})
	if err != nil {
		return domain.JournalEntry{}, err
	}
	if len(entries) == 0 {
		return domain.JournalEntry{}, domain.ErrNotFound
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Week == entries[j].Week {
			return entries[i].UpdatedAt.After(entries[j].UpdatedAt)
		}
		return entries[i].Week > entries[j].Week
	})
	return entries[0], nil
}

func (r *Repository) SaveResources(resources []domain.Resource) error {
	for _, resource := range resources {
		if err := r.SaveResource(resource); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) RemoveResources(entryID string) error {
	resources, err := r.ResourcesFor(entryID)
	if err != nil {
		return err
	}
	for _, resource := range resources {
		if err := r.db.Delete("resources", resource.ID); err != nil {
			return err
		}
	}
	return nil
}
