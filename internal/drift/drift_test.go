package drift_test

import (
	"testing"

	"github.com/yourorg/stackdiff/internal/diff"
	"github.com/yourorg/stackdiff/internal/drift"
)

func makeEntries(added, removed, changed int) []diff.Entry {
	var entries []diff.Entry
	for i := 0; i < added; i++ {
		entries = append(entries, diff.Entry{Key: fmt.Sprintf("add%d", i), Status: diff.StatusAdded, NewVal: "v"})
	}
	for i := 0; i < removed; i++ {
		entries = append(entries, diff.Entry{Key: fmt.Sprintf("rem%d", i), Status: diff.StatusRemoved, OldVal: "v"})
	}
	for i := 0; i < changed; i++ {
		entries = append(entries, diff.Entry{Key: fmt.Sprintf("chg%d", i), Status: diff.StatusChanged, OldVal: "a", NewVal: "b"})
	}
	return entries
}

func TestEvaluate_NoDrift(t *testing.T) {
	r := drift.Evaluate(nil)
	if r.Level != drift.LevelNone {
		t.Errorf("expected none, got %s", r.Level)
	}
	if r.Total != 0 {
		t.Errorf("expected total 0, got %d", r.Total)
	}
}

func TestEvaluate_LowDrift(t *testing.T) {
	entries := makeEntries(1, 0, 1)
	r := drift.Evaluate(entries)
	if r.Level != drift.LevelLow {
		t.Errorf("expected low, got %s", r.Level)
	}
	if r.Total != 2 {
		t.Errorf("expected total 2, got %d", r.Total)
	}
}

func TestEvaluate_ModerateDrift(t *testing.T) {
	entries := makeEntries(2, 1, 2)
	r := drift.Evaluate(entries)
	if r.Level != drift.LevelModerate {
		t.Errorf("expected moderate, got %s", r.Level)
	}
}

func TestEvaluate_HighDrift_ManyTotal(t *testing.T) {
	entries := makeEntries(5, 2, 4)
	r := drift.Evaluate(entries)
	if r.Level != drift.LevelHigh {
		t.Errorf("expected high, got %s", r.Level)
	}
}

func TestEvaluate_HighDrift_ManyRemoved(t *testing.T) {
	entries := makeEntries(0, 3, 0)
	r := drift.Evaluate(entries)
	if r.Level != drift.LevelHigh {
		t.Errorf("expected high due to removals, got %s", r.Level)
	}
}

func TestResult_String(t *testing.T) {
	r := drift.Result{Level: drift.LevelLow, Added: 1, Removed: 0, Changed: 1, Total: 2}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
