package summary_test

import (
	"testing"

	"github.com/yourorg/stackdiff/internal/diff"
	"github.com/yourorg/stackdiff/internal/summary"
)

func entries(statuses ...string) []diff.Entry {
	out := make([]diff.Entry, len(statuses))
	for i, s := range statuses {
		out[i] = diff.Entry{Key: fmt.Sprintf("k%d", i), Status: s}
	}
	return out
}

func TestBuild_NoDrift(t *testing.T) {
	s := summary.Build([]diff.Entry{})
	if s.HasDrift() {
		t.Fatal("expected no drift")
	}
	if s.String() != "no drift detected" {
		t.Fatalf("unexpected string: %s", s.String())
	}
}

func TestBuild_Counts(t *testing.T) {
	entriesSlice := []diff.Entry{
		{Key: "a", Status: diff.StatusAdded},
		{Key: "b", Status: diff.StatusAdded},
		{Key: "c", Status: diff.StatusRemoved},
		{Key: "d", Status: diff.StatusChanged},
		{Key: "e", Status: diff.StatusChanged},
		{Key: "f", Status: diff.StatusChanged},
	}
	s := summary.Build(entriesSlice)
	if s.Added != 2 {
		t.Errorf("expected 2 added, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", s.Removed)
	}
	if s.Changed != 3 {
		t.Errorf("expected 3 changed, got %d", s.Changed)
	}
	if s.Total != 6 {
		t.Errorf("expected total 6, got %d", s.Total)
	}
}

func TestBuild_HasDrift(t *testing.T) {
	s := summary.Build([]diff.Entry{{Key: "x", Status: diff.StatusAdded}})
	if !s.HasDrift() {
		t.Fatal("expected drift")
	}
}

func TestStats_String_ContainsParts(t *testing.T) {
	s := summary.Stats{Added: 1, Removed: 2, Changed: 3, Total: 6}
	str := s.String()
	for _, want := range []string{"1 added", "2 removed", "3 changed", "6"} {
		if !strings.Contains(str, want) {
			t.Errorf("expected %q in %q", want, str)
		}
	}
}
