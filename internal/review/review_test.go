package review

import (
	"path/filepath"
	"testing"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/repository"
	"independentweeklylog/internal/store"
)

func TestReviewDecisionPublishesEntry(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "journal.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	entry := domain.JournalEntry{ID: "reviewed", Week: 20, Title: "Reviewable", Summary: "Review has enough context", Author: "author", Status: domain.StatusSubmitted, Achievements: []string{"x"}, NextSteps: []string{"y"}}
	if err := repo.SaveEntry(entry); err != nil {
		t.Fatal(err)
	}
	updated, _, err := New(repo).Decide(entry.ID, "reviewer", domain.StatusApproved, "looks ready")
	if err != nil || updated.Status != domain.StatusApproved {
		t.Fatalf("review failed: %#v %v", updated, err)
	}
}
