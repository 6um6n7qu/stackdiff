package delta_test

import (
	"strings"
	"testing"

	"github.com/user/stackdiff/internal/delta"
	"github.com/user/stackdiff/internal/diff"
)

func makeEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "A", OldValue: "", NewValue: "1", Status: diff.StatusAdded},
		{Key: "B", OldValue: "x", NewValue: "", Status: diff.StatusRemoved},
		{Key: "C", OldValue: "old", NewValue: "new", Status: diff.StatusChanged},
		{Key: "D", OldValue: "same", NewValue: "same", Status: diff.StatusEqual},
	}
}

func TestCompute_NoDrift(t *testing.T) {
	entries := []diff.Entry{
		{Key: "X", OldValue: "v", NewValue: "v", Status: diff.StatusEqual},
	}
	r := delta.Compute(entries)
	if r.HasDrift() {
		t.Errorf("expected no drift, got %s", r)
	}
}

func TestCompute_Counts(t *testing.T) {
	r := delta.Compute(makeEntries())
	if r.Added != 1 {
		t.Errorf("Added: want 1, got %d", r.Added)
	}
	if r.Removed != 1 {
		t.Errorf("Removed: want 1, got %d", r.Removed)
	}
	if r.Changed != 1 {
		t.Errorf("Changed: want 1, got %d", r.Changed)
	}
	if r.Total != 3 {
		t.Errorf("Total: want 3, got %d", r.Total)
	}
}

func TestCompute_HasDrift(t *testing.T) {
	r := delta.Compute(makeEntries())
	if !r.HasDrift() {
		t.Error("expected HasDrift to be true")
	}
}

func TestResult_String(t *testing.T) {
	r := delta.Compute(makeEntries())
	s := r.String()
	for _, part := range []string{"added=", "removed=", "changed=", "total="} {
		if !strings.Contains(s, part) {
			t.Errorf("String() missing %q, got: %s", part, s)
		}
	}
}
