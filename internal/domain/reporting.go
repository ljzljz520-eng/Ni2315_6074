package domain

import (
	"fmt"
	"sort"
	"strings"
)

type ProgressMetric struct {
	Name   string `json:"name"`
	Value  int    `json:"value"`
	Target int    `json:"target"`
	Unit   string `json:"unit"`
}

type DeveloperReport struct {
	Author   string           `json:"author"`
	Weeks    []int            `json:"weeks"`
	Entries  int              `json:"entries"`
	Approved int              `json:"approved"`
	Archived int              `json:"archived"`
	Metrics  []ProgressMetric `json:"metrics"`
	Themes   []string         `json:"themes"`
}

type StatusBreakdown struct {
	Draft     int `json:"draft"`
	Submitted int `json:"submitted"`
	Approved  int `json:"approved"`
	Rejected  int `json:"rejected"`
	Archived  int `json:"archived"`
}

func CountStatuses(entries []JournalEntry) StatusBreakdown {
	breakdown := StatusBreakdown{}
	for _, entry := range entries {
		switch entry.Status {
		case StatusDraft:
			breakdown.Draft++
		case StatusSubmitted:
			breakdown.Submitted++
		case StatusApproved:
			breakdown.Approved++
		case StatusRejected:
			breakdown.Rejected++
		case StatusArchived:
			breakdown.Archived++
		}
	}
	return breakdown
}

func BuildDeveloperReport(entries []JournalEntry, author string) DeveloperReport {
	report := DeveloperReport{Author: author, Weeks: make([]int, 0), Themes: make([]string, 0)}
	weeks := make(map[int]bool)
	themes := make(map[string]int)
	for _, entry := range entries {
		if author != "" && !strings.EqualFold(entry.Author, author) {
			continue
		}
		report.Entries++
		weeks[entry.Week] = true
		if entry.Status == StatusApproved {
			report.Approved++
		}
		if entry.Status == StatusArchived {
			report.Archived++
		}
		for _, tag := range entry.Tags {
			themes[strings.ToLower(tag)]++
		}
	}
	for week := range weeks {
		report.Weeks = append(report.Weeks, week)
	}
	sort.Ints(report.Weeks)
	type themeCount struct {
		name  string
		count int
	}
	ordered := make([]themeCount, 0, len(themes))
	for name, count := range themes {
		ordered = append(ordered, themeCount{name: name, count: count})
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].count == ordered[j].count {
			return ordered[i].name < ordered[j].name
		}
		return ordered[i].count > ordered[j].count
	})
	for _, item := range ordered {
		report.Themes = append(report.Themes, item.name)
	}
	if report.Entries > 0 {
		report.Metrics = append(report.Metrics, ProgressMetric{Name: "approval-rate", Value: report.Approved * 100 / report.Entries, Target: 80, Unit: "percent"})
		report.Metrics = append(report.Metrics, ProgressMetric{Name: "archive-rate", Value: report.Archived * 100 / report.Entries, Target: 50, Unit: "percent"})
	} else {
		report.Metrics = append(report.Metrics, ProgressMetric{Name: "approval-rate", Target: 80, Unit: "percent"})
	}
	return report
}

func FormatReport(report DeveloperReport) string {
	return fmt.Sprintf("%s: %d entries across %d weeks, %d approved, themes=%s", report.Author, report.Entries, len(report.Weeks), report.Approved, strings.Join(report.Themes, ","))
}

func MergeTags(entries []JournalEntry) []string {
	counts := make(map[string]int)
	for _, entry := range entries {
		for _, tag := range entry.Tags {
			counts[strings.ToLower(strings.TrimSpace(tag))]++
		}
	}
	result := make([]string, 0, len(counts))
	for tag := range counts {
		if tag != "" {
			result = append(result, tag)
		}
	}
	sort.Strings(result)
	return result
}

func EntryText(entry JournalEntry) string {
	parts := []string{entry.Title, entry.Summary, entry.Author}
	parts = append(parts, entry.Achievements...)
	parts = append(parts, entry.Blockers...)
	parts = append(parts, entry.NextSteps...)
	return strings.Join(parts, " ")
}

func IsComplete(entry JournalEntry) bool {
	if strings.TrimSpace(entry.Title) == "" || strings.TrimSpace(entry.Summary) == "" {
		return false
	}
	if len(entry.Achievements) == 0 || len(entry.NextSteps) == 0 {
		return false
	}
	return entry.Status != StatusDraft
}

func CompareProgress(before, after JournalEntry) []string {
	changes := make([]string, 0)
	if before.Status != after.Status {
		changes = append(changes, fmt.Sprintf("status:%s->%s", before.Status, after.Status))
	}
	if before.Summary != after.Summary {
		changes = append(changes, "summary")
	}
	if len(before.Achievements) != len(after.Achievements) {
		changes = append(changes, "achievements")
	}
	if len(before.NextSteps) != len(after.NextSteps) {
		changes = append(changes, "next_steps")
	}
	sort.Strings(changes)
	return changes
}

func ValidateWeekSequence(entries []JournalEntry) error {
	weeks := make([]int, 0, len(entries))
	seen := make(map[int]bool)
	for _, entry := range entries {
		if seen[entry.Week] {
			continue
		}
		seen[entry.Week] = true
		weeks = append(weeks, entry.Week)
	}
	sort.Ints(weeks)
	for i := 1; i < len(weeks); i++ {
		if weeks[i]-weeks[i-1] > 4 {
			return fmt.Errorf("week gap between %d and %d", weeks[i-1], weeks[i])
		}
	}
	return nil
}
