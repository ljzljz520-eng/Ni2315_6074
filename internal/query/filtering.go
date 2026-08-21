package query

import (
	"sort"
	"strings"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/repository"
)

type SortMode string

const (
	SortWeek    SortMode = "week"
	SortUpdated SortMode = "updated"
	SortTitle   SortMode = "title"
)

type Request struct {
	Filter     repository.EntryFilter
	Sort       SortMode
	Descending bool
	Limit      int
}

func (s *Service) Browse(request Request) ([]domain.JournalEntry, error) {
	entries, err := s.repo.FindEntries(request.Filter)
	if err != nil {
		return nil, err
	}
	sortEntries(entries, request.Sort, request.Descending)
	if request.Limit > 0 && len(entries) > request.Limit {
		entries = entries[:request.Limit]
	}
	return entries, nil
}

func sortEntries(entries []domain.JournalEntry, mode SortMode, descending bool) {
	sort.SliceStable(entries, func(i, j int) bool {
		var less bool
		switch mode {
		case SortUpdated:
			less = entries[i].UpdatedAt.Before(entries[j].UpdatedAt)
		case SortTitle:
			less = strings.ToLower(entries[i].Title) < strings.ToLower(entries[j].Title)
		default:
			if entries[i].Week == entries[j].Week {
				less = entries[i].ID < entries[j].ID
			} else {
				less = entries[i].Week < entries[j].Week
			}
		}
		if descending {
			return !less && entries[i].ID != entries[j].ID
		}
		return less
	})
}

func GroupByStatus(entries []domain.JournalEntry) map[domain.EntryStatus][]domain.JournalEntry {
	groups := make(map[domain.EntryStatus][]domain.JournalEntry)
	for _, entry := range entries {
		groups[entry.Status] = append(groups[entry.Status], entry)
	}
	for status := range groups {
		sortEntries(groups[status], SortWeek, false)
	}
	return groups
}

func GroupByTag(entries []domain.JournalEntry) map[string][]domain.JournalEntry {
	groups := make(map[string][]domain.JournalEntry)
	for _, entry := range entries {
		for _, tag := range entry.Tags {
			clean := strings.ToLower(strings.TrimSpace(tag))
			if clean != "" {
				groups[clean] = append(groups[clean], entry)
			}
		}
	}
	return groups
}

func (s *Service) CompletionRate(week int) (float64, error) {
	summary, err := s.Summary(week)
	if err != nil {
		return 0, err
	}
	if summary.Total == 0 {
		return 0, nil
	}
	return float64(summary.Approved+summary.Archived) / float64(summary.Total), nil
}

func (s *Service) Pending(week int) ([]domain.JournalEntry, error) {
	return s.Search(repository.EntryFilter{Week: week, Status: domain.StatusSubmitted})
}

func (s *Service) Visible(week int) ([]domain.JournalEntry, error) {
	entries, err := s.Search(repository.EntryFilter{Week: week})
	if err != nil {
		return nil, err
	}
	result := make([]domain.JournalEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.IsVisible() {
			result = append(result, entry)
		}
	}
	return result, nil
}

func (s *Service) SearchText(text string) ([]domain.JournalEntry, error) {
	return s.Search(repository.EntryFilter{Text: text})
}

func (s *Service) AuthorSummary(author string) (domain.DeveloperReport, error) {
	return s.repo.ReportFor(author)
}
