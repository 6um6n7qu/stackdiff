package report

import (
	"testing"

	"github.com/user/stackdiff/internal/diff"
)

func sampleDiffs() []diff.DiffEntry {
	return []diff.DiffEntry{
		{Key: "APP_ENV", Kind: diff.Added, NewValue: "production"},
		{Key: "DB_HOST", Kind: diff.Removed, OldValue: "localhost"},
		{Key: "LOG_LEVEL", Kind: diff.Changed, OldValue: "debug", NewValue: "info"},
	}
}

func TestNew_Summary(t *testing.T) {
	r := New("staging", "prod", sampleDiffs())
	if r.Summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", r.Summary.Added)
	}
	if r.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", r.Summary.Removed)
	}
	if r.Summary.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", r.Summary.Changed)
	}
	if r.Summary.Total != 3 {
		t.Errorf("expected total 3, got %d", r.Summary.Total)
	}
}

func TestNew_Labels(t *testing.T) {
	r := New("a", "b", nil)
	if r.SourceLabel != "a" || r.TargetLabel != "b" {
		t.Errorf("unexpected labels: %s / %s", r.SourceLabel, r.TargetLabel)
	}
}

func TestHasDrift_True(t *testing.T) {
	r := New("a", "b", sampleDiffs())
	if !r.HasDrift() {
		t.Error("expected HasDrift to be true")
	}
}

func TestHasDrift_False(t *testing.T) {
	r := New("a", "b", []diff.DiffEntry{})
	if r.HasDrift() {
		t.Error("expected HasDrift to be false")
	}
}

func TestNew_EmptyDiffs(t *testing.T) {
	r := New("x", "y", nil)
	if r.Summary.Total != 0 {
		t.Errorf("expected total 0, got %d", r.Summary.Total)
	}
	if r.HasDrift() {
		t.Error("expected no drift for nil diffs")
	}
}
