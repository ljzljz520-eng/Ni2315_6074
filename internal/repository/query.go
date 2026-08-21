package repository

import (
	"strings"

	"independentweeklylog/internal/domain"
)

type EntryFilter struct {
	Week   int
	Author string
	Status domain.EntryStatus
	Tag    string
	Text   string
}

func (r *Repository) FindEntries(filter EntryFilter) ([]domain.JournalEntry, error) {
	entries, err := r.Entries()
	if err != nil {
		return nil, err
	}
	result := make([]domain.JournalEntry, 0, len(entries))
	for _, entry := range entries {
		if filter.Week > 0 && entry.Week != filter.Week {
			continue
		}
		if filter.Author != "" && !strings.EqualFold(entry.Author, filter.Author) {
			continue
		}
		if filter.Status != "" && entry.Status != filter.Status {
			continue
		}
		if filter.Tag != "" && !containsFold(entry.Tags, filter.Tag) {
			continue
		}
		if filter.Text != "" && !textMatch(entry, filter.Text) {
			continue
		}
		result = append(result, entry)
	}
	return result, nil
}

func containsFold(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(value, target) {
			return true
		}
	}
	return false
}

func textMatch(entry domain.JournalEntry, text string) bool {
	needle := strings.ToLower(text)
	if strings.Contains(strings.ToLower(entry.Title), needle) || strings.Contains(strings.ToLower(entry.Summary), needle) {
		return true
	}
	for _, value := range entry.Achievements {
		if strings.Contains(strings.ToLower(value), needle) {
			return true
		}
	}
	return false
}
