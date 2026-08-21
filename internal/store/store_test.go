package store

import (
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "journal.db")
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Put(EntriesBucket, "entry-1", map[string]string{"title": "saved"}); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	var value map[string]string
	if err := second.Get(EntriesBucket, "entry-1", &value); err != nil {
		t.Fatal(err)
	}
	if value["title"] != "saved" {
		t.Fatalf("unexpected value: %#v", value)
	}
}

func TestStoreRejectsUnknownBucket(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "journal.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.Put("unknown", "key", "value"); err == nil {
		t.Fatal("unknown bucket accepted")
	}
}
