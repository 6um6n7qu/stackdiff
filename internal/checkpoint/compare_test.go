package checkpoint_test

import (
	"testing"

	"github.com/user/stackdiff/internal/checkpoint"
	"github.com/user/stackdiff/internal/diff"
)

func makeCheckpoint(entries []diff.Entry) *checkpoint.Checkpoint {
	return checkpoint.New("test", entries)
}

func TestCompare_NoChange(t *testing.T) {
	entries := sampleEntries()
	cp := makeCheckpoint(entries)
	result := checkpoint.Compare(cp, entries)
	if result.HasDrift() {
		t.Errorf("expected no drift, got %d drifted", len(result.Drifted))
	}
}

func TestCompare_AddedEntry(t *testing.T) {
	cp := makeCheckpoint(sampleEntries())
	live := append(sampleEntries(), diff.Entry{Key: "DEBUG", NewValue: "true", Status: diff.StatusEqual})
	result := checkpoint.Compare(cp, live)
	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(result.Added))
	}
	if result.Added[0].Key != "DEBUG" {
		t.Errorf("unexpected added key: %s", result.Added[0].Key)
	}
	if result.Added[0].Status != diff.StatusAdded {
		t.Errorf("expected status added, got %s", result.Added[0].Status)
	}
}

func TestCompare_RemovedEntry(t *testing.T) {
	base := sampleEntries()
	cp := makeCheckpoint(base)
	live := base[:1] // remove PORT
	result := checkpoint.Compare(cp, live)
	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(result.Removed))
	}
	if result.Removed[0].Key != "PORT" {
		t.Errorf("unexpected removed key: %s", result.Removed[0].Key)
	}
}

func TestCompare_ChangedEntry(t *testing.T) {
	cp := makeCheckpoint(sampleEntries())
	live := []diff.Entry{
		{Key: "HOST", OldValue: "localhost", NewValue: "prod.example.com", Status: diff.StatusEqual},
		{Key: "PORT", OldValue: "8080", NewValue: "8080", Status: diff.StatusEqual},
	}
	result := checkpoint.Compare(cp, live)
	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(result.Changed))
	}
	if result.Changed[0].Key != "HOST" {
		t.Errorf("unexpected changed key: %s", result.Changed[0].Key)
	}
	if result.Changed[0].OldValue != "localhost" {
		t.Errorf("expected old value localhost, got %s", result.Changed[0].OldValue)
	}
}

func TestCompare_Summary(t *testing.T) {
	cp := makeCheckpoint(sampleEntries())
	live := []diff.Entry{
		{Key: "HOST", NewValue: "newhost", Status: diff.StatusEqual},
		{Key: "EXTRA", NewValue: "val", Status: diff.StatusEqual},
	}
	result := checkpoint.Compare(cp, live)
	s := result.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}

func TestCompare_HasDrift_False(t *testing.T) {
	cp := makeCheckpoint(nil)
	result := checkpoint.Compare(cp, nil)
	if result.HasDrift() {
		t.Error("expected no drift for empty inputs")
	}
}
