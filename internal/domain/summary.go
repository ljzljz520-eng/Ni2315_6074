package domain

import "sort"

func BuildSummary(entries []JournalEntry, week int) WeeklySummary {
	summary := WeeklySummary{Week: week}
	authors := make(map[string]bool)
	for _, entry := range entries {
		if entry.Week != week {
			continue
		}
		summary.Total++
		authors[entry.Author] = true
		switch entry.Status {
		case StatusApproved:
			summary.Approved++
		case StatusArchived:
			summary.Archived++
		case StatusSubmitted:
			summary.Pending++
		}
	}
	for author := range authors {
		summary.Authors = append(summary.Authors, author)
	}
	sort.Strings(summary.Authors)
	return summary
}
