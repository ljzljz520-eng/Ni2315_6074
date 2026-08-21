package weeklylog_test

import (
	"path/filepath"
	"testing"

	"independentweeklylog/internal/archive"
	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/query"
	"independentweeklylog/internal/repository"
	"independentweeklylog/internal/review"
	"independentweeklylog/internal/store"
	"independentweeklylog/internal/workflow20"
)

func TestWorkflowCaptureReviewQueryArchive(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "journal.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	entry := domain.JournalEntry{ID: "flow", Week: 20, Title: "Playable milestone", Summary: "The weekly build gained a playable loop", Author: "dev", Achievements: []string{"loop"}, Blockers: []string{"audio"}, NextSteps: []string{"polish"}, Tags: []string{"build"}}
	resources := []domain.Resource{{ID: "engine", Name: "engine", Kind: "parent", Payload: "engine"}, {ID: "scene", Name: "scene", Kind: "child", ParentID: "engine", Payload: "scene"}}
	created, err := workflow20.New(repo).CaptureDraft(entry, resources)
	if err != nil {
		t.Fatal(err)
	}
	if created.Entry.Status != domain.StatusDraft {
		t.Fatal("draft not created")
	}
	if _, err := workflow20.New(repo).Submit(entry.ID, "dev"); err != nil {
		t.Fatal(err)
	}
	approved, _, err := review.New(repo).Decide(entry.ID, "producer", domain.StatusApproved, "approved for archive")
	if err != nil || approved.Status != domain.StatusApproved {
		t.Fatalf("approval failed: %#v %v", approved, err)
	}
	items, err := query.New(repo).Search(repository.EntryFilter{Week: 20, Status: domain.StatusApproved})
	if err != nil || len(items) != 1 {
		t.Fatalf("query failed: %#v %v", items, err)
	}
	if _, err := archive.New(repo).Archive(entry.ID, "week complete"); err != nil {
		t.Fatal(err)
	}
	summary, err := query.New(repo).Summary(20)
	if err != nil || summary.Archived != 1 {
		t.Fatalf("archive summary failed: %#v %v", summary, err)
	}
}
