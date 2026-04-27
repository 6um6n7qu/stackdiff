package split_test

import (
	"testing"

	"github.com/stackdiff/stackdiff/internal/diff"
	"github.com/stackdiff/stackdiff/internal/split"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "HOST", OldValue: "a", NewValue: "b", Status: diff.StatusChanged},
		{Key: "PORT", OldValue: "", NewValue: "8080", Status: diff.StatusAdded},
		{Key: "DEBUG", OldValue: "true", NewValue: "", Status: diff.StatusRemoved},
		{Key: "LOG_LEVEL", OldValue: "info", NewValue: "info", Status: diff.StatusEqual},
	}
}

func TestApply_AllMatched(t *testing.T) {
	entries := sampleEntries()
	r := split.Apply("all", entries, func(e diff.Entry) bool { return true })
	if len(r.Matched) != 4 {
		t.Fatalf("expected 4 matched, got %d", len(r.Matched))
	}
	if len(r.Unmatched) != 0 {
		t.Fatalf("expected 0 unmatched, got %d", len(r.Unmatched))
	}
}

func TestApply_NoneMatched(t *testing.T) {
	entries := sampleEntries()
	r := split.Apply("none", entries, func(e diff.Entry) bool { return false })
	if len(r.Matched) != 0 {
		t.Fatalf("expected 0 matched, got %d", len(r.Matched))
	}
	if len(r.Unmatched) != 4 {
		t.Fatalf("expected 4 unmatched, got %d", len(r.Unmatched))
	}
}

func TestApply_PreservesAllEntries(t *testing.T) {
	entries := sampleEntries()
	r := split.Apply("drift", entries, split.ByStatus(diff.StatusChanged, diff.StatusAdded))
	total := len(r.Matched) + len(r.Unmatched)
	if total != len(entries) {
		t.Fatalf("entry count mismatch: got %d, want %d", total, len(entries))
	}
}

func TestByStatus_MatchesCorrectly(t *testing.T) {
	entries := sampleEntries()
	r := split.Apply("changed", entries, split.ByStatus(diff.StatusChanged))
	if len(r.Matched) != 1 {
		t.Fatalf("expected 1 matched, got %d", len(r.Matched))
	}
	if r.Matched[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", r.Matched[0].Key)
	}
}

func TestResult_HasDrift_True(t *testing.T) {
	entries := sampleEntries()
	r := split.Apply("test", entries, split.ByStatus(diff.StatusEqual))
	if !r.HasDrift() {
		t.Error("expected HasDrift true because unmatched contains drift entries")
	}
}

func TestResult_HasDrift_False(t *testing.T) {
	entries := []diff.Entry{
		{Key: "A", OldValue: "x", NewValue: "x", Status: diff.StatusEqual},
	}
	r := split.Apply("test", entries, split.ByStatus(diff.StatusEqual))
	if r.HasDrift() {
		t.Error("expected HasDrift false")
	}
}

func TestMatchedDrift_ReturnsOnlyDrift(t *testing.T) {
	entries := sampleEntries()
	r := split.Apply("all", entries, func(e diff.Entry) bool { return true })
	drift := r.MatchedDrift()
	for _, e := range drift {
		if !e.IsDrift() {
			t.Errorf("non-drift entry in MatchedDrift: %s", e.Key)
		}
	}
	if len(drift) != 3 {
		t.Errorf("expected 3 drift entries, got %d", len(drift))
	}
}

func TestResult_Name(t *testing.T) {
	r := split.Apply("secrets", nil, func(e diff.Entry) bool { return false })
	if r.Name != "secrets" {
		t.Errorf("expected name 'secrets', got %q", r.Name)
	}
}
