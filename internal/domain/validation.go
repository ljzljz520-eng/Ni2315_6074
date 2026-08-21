package domain

import (
	"fmt"
	"sort"
	"strings"
)

type ValidationIssue struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type EntryValidator struct {
	MinTags      int
	MaxNextSteps int
}

func NewEntryValidator() EntryValidator {
	return EntryValidator{MinTags: 1, MaxNextSteps: 8}
}

func (v EntryValidator) Check(entry JournalEntry) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	if strings.TrimSpace(entry.Title) == "" {
		issues = append(issues, ValidationIssue{Field: "title", Message: "title is required", Severity: "error"})
	}
	if len(strings.TrimSpace(entry.Summary)) < 12 {
		issues = append(issues, ValidationIssue{Field: "summary", Message: "summary must be at least 12 characters", Severity: "error"})
	}
	if len(entry.Achievements) == 0 {
		issues = append(issues, ValidationIssue{Field: "achievements", Message: "record an achievement", Severity: "error"})
	}
	if len(entry.NextSteps) == 0 {
		issues = append(issues, ValidationIssue{Field: "next_steps", Message: "record a next step", Severity: "error"})
	}
	if len(entry.NextSteps) > v.MaxNextSteps {
		issues = append(issues, ValidationIssue{Field: "next_steps", Message: fmt.Sprintf("at most %d next steps", v.MaxNextSteps), Severity: "warning"})
	}
	if len(entry.Tags) < v.MinTags {
		issues = append(issues, ValidationIssue{Field: "tags", Message: "add a tag for discovery", Severity: "warning"})
	}
	if len(entry.Blockers) > 6 {
		issues = append(issues, ValidationIssue{Field: "blockers", Message: "summarize blockers before review", Severity: "warning"})
	}
	return issues
}

func HasErrors(issues []ValidationIssue) bool {
	for _, issue := range issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

func NormalizeList(values []string) []string {
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		clean := strings.TrimSpace(value)
		if clean == "" || seen[strings.ToLower(clean)] {
			continue
		}
		seen[strings.ToLower(clean)] = true
		result = append(result, clean)
	}
	sort.Strings(result)
	return result
}

func ValidateResourceTree(resources []Resource) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	ids := make(map[string]bool, len(resources))
	for _, resource := range resources {
		if ids[resource.ID] {
			issues = append(issues, ValidationIssue{Field: resource.ID, Message: "duplicate resource id", Severity: "error"})
		}
		ids[resource.ID] = true
	}
	for _, resource := range resources {
		if resource.ParentID != "" && !ids[resource.ParentID] {
			issues = append(issues, ValidationIssue{Field: resource.ID, Message: "parent resource is missing", Severity: "error"})
		}
	}
	return issues
}
