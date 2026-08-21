package query

import (
	"path/filepath"
	"testing"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/repository"
	"independentweeklylog/internal/store"
)

func TestSummaryCountsVisibleStates(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "journal.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	for i, status := range []domain.EntryStatus{domain.StatusApproved, domain.StatusSubmitted, domain.StatusArchived} {
		entry := domain.JournalEntry{ID: string(rune('a' + i)), Week: 20, Title: "Entry", Summary: "Summary with enough words", Author: string(rune('a' + i)), Status: status, Achievements: []string{"x"}, NextSteps: []string{"y"}}
		if err := repo.SaveEntry(entry); err != nil {
			t.Fatal(err)
		}
	}
	summary, err := New(repo).Summary(20)
	if err != nil || summary.Total != 3 || summary.Approved != 1 || summary.Pending != 1 || summary.Archived != 1 {
		t.Fatalf("unexpected summary: %#v %v", summary, err)
	}
}
