package history_test

import (
	"os"
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/history"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "LOG_LEVEL", Left: "debug", Right: "info", Status: diff.StatusChanged},
		{Key: "PORT", Left: "8080", Right: "8080", Status: diff.StatusEqual},
		{Key: "NEW_KEY", Left: "", Right: "value", Status: diff.StatusAdded},
	}
}

func TestNewRecord_Fields(t *testing.T) {
	r := history.NewRecord("staging", "prod", sampleEntries())
	if r.LeftLabel != "staging" {
		t.Errorf("expected LeftLabel=staging, got %s", r.LeftLabel)
	}
	if r.RightLabel != "prod" {
		t.Errorf("expected RightLabel=prod, got %s", r.RightLabel)
	}
	if r.ID == "" {
		t.Error("expected non-empty ID")
	}
	if len(r.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(r.Entries))
	}
}

func TestRecord_DriftCount(t *testing.T) {
	r := history.NewRecord("a", "b", sampleEntries())
	if got := r.DriftCount(); got != 2 {
		t.Errorf("expected DriftCount=2, got %d", got)
	}
}

func TestRecord_HasDrift(t *testing.T) {
	r := history.NewRecord("a", "b", sampleEntries())
	if !r.HasDrift() {
		t.Error("expected HasDrift=true")
	}
}

func TestRecord_Summary(t *testing.T) {
	r := history.NewRecord("staging", "prod", sampleEntries())
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := history.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	r := history.NewRecord("left", "right", sampleEntries())
	if err := store.Save(r); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := store.Load(r.ID)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.ID != r.ID {
		t.Errorf("ID mismatch: got %s, want %s", loaded.ID, r.ID)
	}
	if loaded.LeftLabel != r.LeftLabel {
		t.Errorf("LeftLabel mismatch: got %s", loaded.LeftLabel)
	}
}

func TestStore_List(t *testing.T) {
	dir := t.TempDir()
	store, _ := history.NewStore(dir)
	for i := 0; i < 3; i++ {
		r := history.NewRecord("a", "b", sampleEntries())
		_ = store.Save(r)
	}
	ids, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(ids) != 3 {
		t.Errorf("expected 3 records, got %d", len(ids))
	}
}

func TestStore_Load_Missing(t *testing.T) {
	dir := t.TempDir()
	store, _ := history.NewStore(dir)
	_, err := store.Load("nonexistent")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}
