package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/stackdiff/internal/baseline"
	"github.com/user/stackdiff/internal/diff"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "LOG_LEVEL", OldValue: "", NewValue: "info", Status: diff.StatusAdded},
		{Key: "PORT", OldValue: "", NewValue: "8080", Status: diff.StatusAdded},
	}
}

func TestNew_Fields(t *testing.T) {
	entries := sampleEntries()
	b := baseline.New("prod", entries)
	if b.Name != "prod" {
		t.Errorf("expected name prod, got %s", b.Name)
	}
	if len(b.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(b.Entries))
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	b := baseline.New("staging", sampleEntries())
	b.CreatedAt = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	if err := baseline.Save(b, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Name != b.Name {
		t.Errorf("name mismatch: got %s", loaded.Name)
	}
	if len(loaded.Entries) != len(b.Entries) {
		t.Errorf("entries mismatch: got %d", len(loaded.Entries))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	b := baseline.New("x", nil)
	err := baseline.Save(b, "/nonexistent/dir/baseline.json")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestDriftFrom_NoDrift(t *testing.T) {
	b := baseline.New("prod", sampleEntries())
	drifted := b.DriftFrom(sampleEntries())
	if len(drifted) != 0 {
		t.Errorf("expected no drift, got %d entries", len(drifted))
	}
}

func TestDriftFrom_DetectsDrift(t *testing.T) {
	b := baseline.New("prod", sampleEntries())
	current := []diff.Entry{
		{Key: "LOG_LEVEL", OldValue: "", NewValue: "debug", Status: diff.StatusChanged},
		{Key: "PORT", OldValue: "", NewValue: "8080", Status: diff.StatusAdded},
	}
	drifted := b.DriftFrom(current)
	if len(drifted) != 1 {
		t.Fatalf("expected 1 drifted entry, got %d", len(drifted))
	}
	if drifted[0].Key != "LOG_LEVEL" {
		t.Errorf("expected LOG_LEVEL, got %s", drifted[0].Key)
	}
	_ = os.Getenv("CI") // suppress unused import
}
