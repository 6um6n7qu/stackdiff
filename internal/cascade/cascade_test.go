package cascade_test

import (
	"errors"
	"testing"

	"github.com/yourusername/stackdiff/internal/cascade"
	"github.com/yourusername/stackdiff/internal/diff"
)

func makeEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "A", OldValue: "1", NewValue: "2", Status: diff.StatusChanged},
		{Key: "B", OldValue: "", NewValue: "3", Status: diff.StatusAdded},
		{Key: "C", OldValue: "x", NewValue: "x", Status: diff.StatusEqual},
	}
}

func TestRun_NoStages_ReturnsSameEntries(t *testing.T) {
	entries := makeEntries()
	report, err := cascade.Run(entries, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Final) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(report.Final))
	}
	if len(report.Stages) != 0 {
		t.Errorf("expected no stage results, got %d", len(report.Stages))
	}
}

func TestRun_SingleStage_FiltersEntries(t *testing.T) {
	onlyDrift := cascade.Stage{
		Name: "only-drift",
		Fn: func(in []diff.Entry) ([]diff.Entry, error) {
			var out []diff.Entry
			for _, e := range in {
				if e.IsDrift() {
					out = append(out, e)
				}
			}
			return out, nil
		},
	}

	report, err := cascade.Run(makeEntries(), []cascade.Stage{onlyDrift})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Final) != 2 {
		t.Errorf("expected 2 drift entries, got %d", len(report.Final))
	}
	if report.Stages[0].Dropped != 1 {
		t.Errorf("expected 1 dropped entry, got %d", report.Stages[0].Dropped)
	}
}

func TestRun_MultipleStages_ChainedCorrectly(t *testing.T) {
	count := 0
	s1 := cascade.Stage{Name: "s1", Fn: func(in []diff.Entry) ([]diff.Entry, error) { count++; return in, nil }}
	s2 := cascade.Stage{Name: "s2", Fn: func(in []diff.Entry) ([]diff.Entry, error) { count++; return in, nil }}

	_, err := cascade.Run(makeEntries(), []cascade.Stage{s1, s2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected both stages to run, count=%d", count)
	}
}

func TestRun_StageError_AbortsEarly(t *testing.T) {
	fail := cascade.Stage{
		Name: "fail",
		Fn:   func(_ []diff.Entry) ([]diff.Entry, error) { return nil, errors.New("boom") },
	}
	_, err := cascade.Run(makeEntries(), []cascade.Stage{fail})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRun_NilStageFn_ReturnsError(t *testing.T) {
	_, err := cascade.Run(makeEntries(), []cascade.Stage{{Name: "nil-fn", Fn: nil}})
	if err == nil {
		t.Fatal("expected error for nil stage function")
	}
}

func TestReport_HasDrift_True(t *testing.T) {
	report, _ := cascade.Run(makeEntries(), nil)
	if !report.HasDrift() {
		t.Error("expected HasDrift to be true")
	}
}

func TestReport_HasDrift_False(t *testing.T) {
	equal := []diff.Entry{
		{Key: "X", OldValue: "1", NewValue: "1", Status: diff.StatusEqual},
	}
	report, _ := cascade.Run(equal, nil)
	if report.HasDrift() {
		t.Error("expected HasDrift to be false")
	}
}
