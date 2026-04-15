package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stackdiff/internal/diff"
	"github.com/stackdiff/internal/snapshot"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "PORT", ValueA: "8080", ValueB: "8080", Status: diff.StatusEqual},
		{Key: "DEBUG", ValueA: "true", ValueB: "false", Status: diff.StatusChanged},
	}
}

func TestNew_Fields(t *testing.T) {
	s := snapshot.New("staging", "prod", sampleEntries())
	if s.LabelA != "staging" {
		t.Errorf("expected LabelA=staging, got %s", s.LabelA)
	}
	if s.LabelB != "prod" {
		t.Errorf("expected LabelB=prod, got %s", s.LabelB)
	}
	if len(s.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(s.Entries))
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.New("a", "b", sampleEntries())
	orig.Timestamp = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	if err := snapshot.Save(orig, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.LabelA != orig.LabelA || loaded.LabelB != orig.LabelB {
		t.Errorf("labels mismatch: got %s/%s", loaded.LabelA, loaded.LabelB)
	}
	if len(loaded.Entries) != len(orig.Entries) {
		t.Errorf("entries count mismatch: got %d", len(loaded.Entries))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	s := snapshot.New("a", "b", nil)
	err := snapshot.Save(s, "/nonexistent/dir/snap.json")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0644)
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
