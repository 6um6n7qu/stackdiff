package reorder_test

import (
	"testing"

	"github.com/stackdiff/stackdiff/internal/diff"
	"github.com/stackdiff/stackdiff/internal/reorder"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "zoo", Status: diff.StatusEqual},
		{Key: "alpha", Status: diff.StatusAdded},
		{Key: "beta", Status: diff.StatusRemoved},
		{Key: "gamma", Status: diff.StatusChanged},
		{Key: "delta", Status: diff.StatusEqual},
	}
}

func TestApply_ByKey_Ascending(t *testing.T) {
	entries := sampleEntries()
	cfg := reorder.Config{Strategy: reorder.ByKey}
	out := reorder.Apply(entries, cfg)

	keys := make([]string, len(out))
	for i, e := range out {
		keys[i] = e.Key
	}
	expected := []string{"alpha", "beta", "delta", "gamma", "zoo"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: got %q, want %q", i, keys[i], k)
		}
	}
}

func TestApply_ByKeyDesc(t *testing.T) {
	entries := sampleEntries()
	cfg := reorder.Config{Strategy: reorder.ByKeyDesc}
	out := reorder.Apply(entries, cfg)

	if out[0].Key != "zoo" {
		t.Errorf("expected first key to be 'zoo', got %q", out[0].Key)
	}
	if out[len(out)-1].Key != "alpha" {
		t.Errorf("expected last key to be 'alpha', got %q", out[len(out)-1].Key)
	}
}

func TestApply_ByStatus_ChangedFirst(t *testing.T) {
	entries := sampleEntries()
	cfg := reorder.DefaultConfig()
	out := reorder.Apply(entries, cfg)

	if out[0].Status != diff.StatusChanged {
		t.Errorf("expected first entry to be StatusChanged, got %v", out[0].Status)
	}
}

func TestApply_ByStatus_EqualLast(t *testing.T) {
	entries := sampleEntries()
	cfg := reorder.DefaultConfig()
	out := reorder.Apply(entries, cfg)

	n := len(out)
	if out[n-1].Status != diff.StatusEqual || out[n-2].Status != diff.StatusEqual {
		t.Errorf("expected last two entries to be StatusEqual")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	entries := sampleEntries()
	origFirst := entries[0].Key
	cfg := reorder.Config{Strategy: reorder.ByKey}
	reorder.Apply(entries, cfg)

	if entries[0].Key != origFirst {
		t.Errorf("Apply mutated the original slice")
	}
}

func TestApply_EmptySlice(t *testing.T) {
	out := reorder.Apply([]diff.Entry{}, reorder.DefaultConfig())
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(out))
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := reorder.DefaultConfig()
	if cfg.Strategy != reorder.ByStatus {
		t.Errorf("expected default strategy ByStatus, got %v", cfg.Strategy)
	}
	if !cfg.StableEqual {
		t.Error("expected StableEqual to be true by default")
	}
}
