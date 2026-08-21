package archive

import (
	"path/filepath"
	"testing"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/repository"
	"independentweeklylog/internal/store"
)

func TestArchiveStoresChecksum(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "journal.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	entry := domain.JournalEntry{ID: "archived", Week: 20, Title: "Archive me", Summary: "Stable snapshot for the archive", Author: "dev", Status: domain.StatusApproved, Achievements: []string{"x"}, NextSteps: []string{"y"}}
	if err := repo.SaveEntry(entry); err != nil {
		t.Fatal(err)
	}
	for _, resource := range []domain.Resource{{ID: "parent", EntryID: entry.ID, Name: "engine", Payload: "a"}, {ID: "child", EntryID: entry.ID, Name: "scene", ParentID: "parent", Payload: "b"}} {
		if err := repo.SaveResource(resource); err != nil {
			t.Fatal(err)
		}
	}
	record, err := New(repo).Archive(entry.ID, "week complete")
	if err != nil || record.Checksum == "" {
		t.Fatalf("archive failed: %#v %v", record, err)
	}
	updated, err := repo.Entry(entry.ID)
	if err != nil || updated.Status != domain.StatusArchived {
		t.Fatalf("entry not archived: %#v %v", updated, err)
	}
}
