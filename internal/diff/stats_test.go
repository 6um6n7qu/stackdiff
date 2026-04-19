package diff

import (
	"testing"
)

func makeStatsEntries() []Entry {
	return []Entry{
		{Key: "a", OldVal: "", NewVal: "1", Status: StatusAdded},
		{Key: "b", OldVal: "2", NewVal: "", Status: StatusRemoved},
		{Key: "c", OldVal: "3", NewVal: "4", Status: StatusChanged},
		{Key: "d", OldVal: "5", NewVal: "5", Status: StatusEqual},
		{Key: "e", OldVal: "6", NewVal: "6", Status: StatusEqual},
	}
}

func TestComputeStats_Counts(t *testing.T) {
	s := ComputeStats(makeStatsEntries())
	if s.Added != 1 {
		t.Errorf("Added: want 1, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("Removed: want 1, got %d", s.Removed)
	}
	if s.Changed != 1 {
		t.Errorf("Changed: want 1, got %d", s.Changed)
	}
	if s.Equal != 2 {
		t.Errorf("Equal: want 2, got %d", s.Equal)
	}
	if s.Total != 5 {
		t.Errorf("Total: want 5, got %d", s.Total)
	}
}

func TestComputeStats_HasDrift(t *testing.T) {
	s := ComputeStats(makeStatsEntries())
	if !s.HasDrift() {
		t.Error("expected HasDrift to be true")
	}
}

func TestComputeStats_NoDrift(t *testing.T) {
	entries := []Entry{
		{Key: "x", OldVal: "1", NewVal: "1", Status: StatusEqual},
	}
	s := ComputeStats(entries)
	if s.HasDrift() {
		t.Error("expected HasDrift to be false")
	}
	if s.DriftCount() != 0 {
		t.Errorf("DriftCount: want 0, got %d", s.DriftCount())
	}
}

func TestComputeStats_DriftCount(t *testing.T) {
	s := ComputeStats(makeStatsEntries())
	if s.DriftCount() != 3 {
		t.Errorf("DriftCount: want 3, got %d", s.DriftCount())
	}
}

func TestComputeStats_Empty(t *testing.T) {
	s := ComputeStats([]Entry{})
	if s.Total != 0 || s.HasDrift() {
		t.Error("expected empty stats with no drift")
	}
}
