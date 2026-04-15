package snapshot_test

import (
	"testing"

	"github.com/stackdiff/internal/diff"
	"github.com/stackdiff/internal/snapshot"
)

func makeSnapshot(label string, entries []diff.Entry) *snapshot.Snapshot {
	return snapshot.New(label, label+"-b", entries)
}

func TestCompareSnapshots_NoChange(t *testing.T) {
	entries := sampleEntries()
	old := makeSnapshot("v1", entries)
	newer := makeSnapshot("v2", entries)

	delta, err := snapshot.CompareSnapshots(old, newer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(delta.Deltas) != 0 {
		t.Errorf("expected 0 deltas, got %d", len(delta.Deltas))
	}
}

func TestCompareSnapshots_AddedEntry(t *testing.T) {
	old := makeSnapshot("v1", sampleEntries())
	newEntries := append(sampleEntries(), diff.Entry{Key: "NEW_KEY", ValueA: "x", ValueB: "x", Status: diff.StatusEqual})
	newer := makeSnapshot("v2", newEntries)

	delta, err := snapshot.CompareSnapshots(old, newer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(delta.Deltas) != 1 || delta.Deltas[0].Change != "added" {
		t.Errorf("expected 1 added delta, got %+v", delta.Deltas)
	}
}

func TestCompareSnapshots_RemovedEntry(t *testing.T) {
	old := makeSnapshot("v1", sampleEntries())
	newer := makeSnapshot("v2", sampleEntries()[1:])

	delta, err := snapshot.CompareSnapshots(old, newer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(delta.Deltas) != 1 || delta.Deltas[0].Change != "removed" {
		t.Errorf("expected 1 removed delta, got %+v", delta.Deltas)
	}
}

func TestCompareSnapshots_ModifiedEntry(t *testing.T) {
	old := makeSnapshot("v1", sampleEntries())
	modified := []diff.Entry{
		{Key: "PORT", ValueA: "8080", ValueB: "9090", Status: diff.StatusChanged},
		{Key: "DEBUG", ValueA: "true", ValueB: "false", Status: diff.StatusChanged},
	}
	newer := makeSnapshot("v2", modified)

	delta, err := snapshot.CompareSnapshots(old, newer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(delta.Deltas) != 1 || delta.Deltas[0].Change != "modified" {
		t.Errorf("expected 1 modified delta, got %+v", delta.Deltas)
	}
}

func TestCompareSnapshots_NilInputs(t *testing.T) {
	_, err := snapshot.CompareSnapshots(nil, nil)
	if err == nil {
		t.Error("expected error for nil snapshots")
	}
}
