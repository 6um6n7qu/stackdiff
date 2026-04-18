package score_test

import (
	"strings"
	"testing"

	"github.com/yourorg/stackdiff/internal/diff"
	"github.com/yourorg/stackdiff/internal/score"
)

func makeEntries(statuses ...string) []diff.Entry {
	var entries []diff.Entry
	for i, s := range statuses {
		entries = append(entries, diff.Entry{
			Key:    fmt.Sprintf("KEY_%d", i),
			Status: s,
		})
	}
	return entries
}

func TestCompute_NoDrift(t *testing.T) {
	r := score.Compute(nil)
	if r.Total != 0 {
		t.Errorf("expected 0, got %d", r.Total)
	}
	if r.Grade() != "A" {
		t.Errorf("expected grade A, got %s", r.Grade())
	}
}

func TestCompute_ChangedWeightsDouble(t *testing.T) {
	entries := []diff.Entry{
		{Key: "K", Status: diff.StatusChanged},
	}
	r := score.Compute(entries)
	if r.Total != 2 {
		t.Errorf("expected 2, got %d", r.Total)
	}
	if r.Changed != 1 {
		t.Errorf("expected Changed=1, got %d", r.Changed)
	}
}

func TestCompute_AddedAndRemoved(t *testing.T) {
	entries := []diff.Entry{
		{Key: "A", Status: diff.StatusAdded},
		{Key: "B", Status: diff.StatusRemoved},
	}
	r := score.Compute(entries)
	if r.Total != 2 {
		t.Errorf("expected 2, got %d", r.Total)
	}
	if r.Added != 1 || r.Removed != 1 {
		t.Errorf("unexpected counts: added=%d removed=%d", r.Added, r.Removed)
	}
}

func TestGrade_Boundaries(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{0, "A"}, {3, "B"}, {8, "C"}, {15, "D"}, {16, "F"},
	}
	for _, tc := range cases {
		r := score.Result{Total: tc.score}
		if g := r.Grade(); g != tc.want {
			t.Errorf("score %d: expected %s got %s", tc.score, tc.want, g)
		}
	}
}

func TestResult_String(t *testing.T) {
	r := score.Result{Total: 5, Changed: 2, Added: 1, Removed: 0}
	s := r.String()
	if !strings.Contains(s, "score=5") {
		t.Errorf("missing score in: %s", s)
	}
	if !strings.Contains(s, "grade=") {
		t.Errorf("missing grade in: %s", s)
	}
}
