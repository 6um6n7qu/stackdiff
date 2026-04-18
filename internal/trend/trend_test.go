package trend_test

import (
	"testing"
	"time"

	"github.com/stackdiff/internal/diff"
	"github.com/stackdiff/internal/history"
	"github.com/stackdiff/internal/trend"
)

func makeRecord(ts time.Time, keys []string, status diff.Status) history.Record {
	var diffs []diff.Entry
	for _, k := range keys {
		diffs = append(diffs, diff.Entry{Key: k, OldVal: "a", NewVal: "b", Status: status})
	}
	return history.Record{Timestamp: ts, Diffs: diffs}
}

func TestAnalyze_NoDrift(t *testing.T) {
	r := trend.Analyze(nil, 24*time.Hour)
	if len(r.Entries) != 0 {
		t.Fatalf("expected no entries, got %d", len(r.Entries))
	}
}

func TestAnalyze_SingleOccurrenceExcluded(t *testing.T) {
	now := time.Now()
	records := []history.Record{
		makeRecord(now.Add(-1*time.Hour), []string{"DB_HOST"}, diff.StatusChanged),
	}
	r := trend.Analyze(records, 24*time.Hour)
	if len(r.Entries) != 0 {
		t.Fatalf("expected single-occurrence key to be excluded")
	}
}

func TestAnalyze_RepeatedKeyIncluded(t *testing.T) {
	now := time.Now()
	records := []history.Record{
		makeRecord(now.Add(-2*time.Hour), []string{"DB_HOST"}, diff.StatusChanged),
		makeRecord(now.Add(-1*time.Hour), []string{"DB_HOST"}, diff.StatusChanged),
	}
	r := trend.Analyze(records, 24*time.Hour)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].Key != "DB_HOST" {
		t.Errorf("unexpected key %s", r.Entries[0].Key)
	}
	if r.Entries[0].Count != 2 {
		t.Errorf("expected count 2, got %d", r.Entries[0].Count)
	}
}

func TestAnalyze_OutsideWindowExcluded(t *testing.T) {
	now := time.Now()
	records := []history.Record{
		makeRecord(now.Add(-48*time.Hour), []string{"API_KEY"}, diff.StatusChanged),
		makeRecord(now.Add(-47*time.Hour), []string{"API_KEY"}, diff.StatusChanged),
	}
	r := trend.Analyze(records, 24*time.Hour)
	if len(r.Entries) != 0 {
		t.Fatalf("expected records outside window to be excluded")
	}
}

func TestAnalyze_SortedByCount(t *testing.T) {
	now := time.Now()
	records := []history.Record{
		makeRecord(now.Add(-3*time.Hour), []string{"X", "Y"}, diff.StatusChanged),
		makeRecord(now.Add(-2*time.Hour), []string{"X", "Y"}, diff.StatusChanged),
		makeRecord(now.Add(-1*time.Hour), []string{"X"}, diff.StatusChanged),
	}
	r := trend.Analyze(records, 24*time.Hour)
	if len(r.Entries) < 2 {
		t.Fatalf("expected at least 2 entries")
	}
	if r.Entries[0].Key != "X" {
		t.Errorf("expected X to rank first, got %s", r.Entries[0].Key)
	}
}
