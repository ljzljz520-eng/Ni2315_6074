package query

import (
	"sort"

	"independentweeklylog/internal/domain"
)

type CalendarWeek struct {
	Week     int                   `json:"week"`
	Entries  []domain.JournalEntry `json:"entries"`
	Complete bool                  `json:"complete"`
}

func BuildCalendar(entries []domain.JournalEntry) []CalendarWeek {
	groups := make(map[int][]domain.JournalEntry)
	for _, entry := range entries {
		groups[entry.Week] = append(groups[entry.Week], entry)
	}
	weeks := make([]int, 0, len(groups))
	for week := range groups {
		weeks = append(weeks, week)
	}
	sort.Ints(weeks)
	calendar := make([]CalendarWeek, 0, len(weeks))
	for _, week := range weeks {
		items := groups[week]
		sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
		complete := len(items) > 0
		for _, item := range items {
			if !domain.IsComplete(item) {
				complete = false
			}
		}
		calendar = append(calendar, CalendarWeek{Week: week, Entries: items, Complete: complete})
	}
	return calendar
}

func MissingWeeks(entries []domain.JournalEntry, start, end int) []int {
	seen := make(map[int]bool)
	for _, entry := range entries {
		seen[entry.Week] = true
	}
	missing := make([]int, 0)
	for week := start; week <= end; week++ {
		if !seen[week] {
			missing = append(missing, week)
		}
	}
	return missing
}
