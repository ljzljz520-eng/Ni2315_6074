package workflow20

import (
	"path/filepath"
	"testing"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/repository"
	"independentweeklylog/internal/store"
)

func TestWorkflow20BusinessInvariant(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "journal.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	entry := domain.JournalEntry{ID: "week20", Week: 20, Title: "Resource orchestration", Summary: "Parent and child resources are saved", Author: "dev", Status: domain.StatusDraft, Achievements: []string{"resource"}, NextSteps: []string{"review"}}
	if err := repo.SaveEntry(entry); err != nil {
		t.Fatal(err)
	}
	resources := []domain.Resource{{ID: "parent", EntryID: entry.ID, Name: "engine", Kind: "parent", Payload: "engine-state"}, {ID: "child", EntryID: entry.ID, Name: "scene", Kind: "child", ParentID: "parent", Payload: "scene-state"}}
	log := &CloseLog{}
	result, err := New(repo).SaveResourceStep(entry.ID, resources, log)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.CloseOrder) > 0 {
		t.Fatal("unexpected close order")
	}
	saved, err := repo.Resource("child")
	if err != nil {
		t.Fatal(err)
	}
	if saved.Payload != "scene-state" {
		t.Fatalf("child resource payload was not preserved: %q", saved.Payload)
	}
	if !log.ValidParentFirst("parent", "child") {
		t.Fatalf("resources closed in wrong order: %#v", log.Items)
	}
}

func TestWorkflowStepsRequireOrder(t *testing.T) {
	steps := DefaultSteps()
	if _, err := CompleteStep(steps, 7); err == nil {
		t.Fatal("out of order step accepted")
	}
	steps, err := CompleteStep(steps, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !steps[0].Complete {
		t.Fatal("step was not completed")
	}
}
