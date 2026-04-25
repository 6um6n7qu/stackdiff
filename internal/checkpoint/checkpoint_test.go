package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/stackdiff/internal/checkpoint"
	"github.com/user/stackdiff/internal/diff"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "HOST", OldValue: "localhost", NewValue: "localhost", Status: diff.StatusEqual},
		{Key: "PORT", OldValue: "8080", NewValue: "8080", Status: diff.StatusEqual},
	}
}

func TestNew_Fields(t *testing.T) {
	entries := sampleEntries()
	cp := checkpoint.New("v1", entries)
	if cp.Name != "v1" {
		t.Fatalf("expected name v1, got %s", cp.Name)
	}
	if len(cp.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(cp.Entries))
	}
	if cp.CreatedAt.IsZero() {
		t.Fatal("expected non-zero CreatedAt")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	cp := checkpoint.New("deploy", sampleEntries())
	cp.CreatedAt = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	if err := checkpoint.Save(dir, "deploy", cp); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := checkpoint.Load(dir, "deploy")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Name != cp.Name {
		t.Errorf("name mismatch: got %s", loaded.Name)
	}
	if len(loaded.Entries) != len(cp.Entries) {
		t.Errorf("entries mismatch: got %d", len(loaded.Entries))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := checkpoint.Load(t.TempDir(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	cp := checkpoint.New("x", nil)
	err := checkpoint.Save("/dev/null/bad", "x", cp)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestList_ReturnsNames(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"alpha", "beta", "gamma"} {
		cp := checkpoint.New(name, sampleEntries())
		if err := checkpoint.Save(dir, name, cp); err != nil {
			t.Fatalf("Save %s: %v", name, err)
		}
	}
	names, err := checkpoint.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 names, got %d", len(names))
	}
}

func TestList_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	names, err := checkpoint.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestList_MissingDir(t *testing.T) {
	names, err := checkpoint.List(filepath.Join(os.TempDir(), "no-such-dir-xyz"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if names != nil {
		t.Errorf("expected nil, got %v", names)
	}
}
