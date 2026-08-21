package repository

import (
	"path/filepath"
	"testing"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/store"
)

func TestRepositoryRoundTrip(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "journal.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := New(db)
	entry := domain.JournalEntry{ID: "entry-1", Week: 20, Title: "Playable slice", Summary: "A useful weekly summary", Author: "dev", Achievements: []string{"build"}, NextSteps: []string{"polish"}}
	if err := repo.SaveEntry(entry); err != nil {
		t.Fatal(err)
	}
	loaded, err := repo.Entry(entry.ID)
	if err != nil || loaded.Title != entry.Title {
		t.Fatalf("round trip failed: %#v %v", loaded, err)
	}
}

func TestRepositoryFiltersEntries(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "journal.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := New(db)
	for _, entry := range []domain.JournalEntry{{ID: "a", Week: 20, Title: "Combat", Summary: "Combat progress summary", Author: "dev", Tags: []string{"combat"}, Achievements: []string{"x"}, NextSteps: []string{"y"}}, {ID: "b", Week: 21, Title: "UI", Summary: "Interface progress summary", Author: "dev", Tags: []string{"ui"}, Achievements: []string{"x"}, NextSteps: []string{"y"}}} {
		if err := repo.SaveEntry(entry); err != nil {
			t.Fatal(err)
		}
	}
	items, err := repo.FindEntries(EntryFilter{Week: 20, Tag: "combat"})
	if err != nil || len(items) != 1 || items[0].ID != "a" {
		t.Fatalf("unexpected filter result: %#v %v", items, err)
	}
}
