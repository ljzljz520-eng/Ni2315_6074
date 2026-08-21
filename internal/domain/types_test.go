package domain

import "testing"

func TestJournalEntryValidation(t *testing.T) {
	entry := JournalEntry{ID: "e1", Week: 20, Title: "Week twenty", Summary: "Implemented playable build notes", Author: "dev", Achievements: []string{"combat"}, NextSteps: []string{"test"}}
	if err := entry.Validate(); err != nil {
		t.Fatalf("valid entry rejected: %v", err)
	}
	entry.Summary = "short"
	if err := entry.Validate(); err == nil {
		t.Fatal("short summary accepted")
	}
}

func TestNormalizeListAndSummary(t *testing.T) {
	values := NormalizeList([]string{" combat ", "Combat", "ui"})
	if len(values) != 2 || values[0] != "combat" || values[1] != "ui" {
		t.Fatalf("unexpected normalized values: %#v", values)
	}
	entries := []JournalEntry{{Week: 20, Author: "a", Status: StatusApproved}, {Week: 20, Author: "b", Status: StatusSubmitted}}
	summary := BuildSummary(entries, 20)
	if summary.Total != 2 || summary.Approved != 1 || summary.Pending != 1 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
}
