package stale_test

import (
	"testing"
	"time"

	"stackdiff/internal/diff"
	"stackdiff/internal/stale"
)

func makeEntry(key, oldVal, newVal string, status diff.Status) diff.Entry {
	return diff.Entry{Key: key, OldValue: oldVal, NewValue: newVal, Status: status}
}

func TestDetect_NoDrift(t *testing.T) {
	snaps := [][]diff.Entry{
		{{Key: "A", Status: diff.StatusEqual}},
		{{Key: "A", Status: diff.StatusEqual}},
		{{Key: "A", Status: diff.StatusEqual}},
	}
	ts := []time.Time{
		time.Now().Add(-100 * time.Hour),
		time.Now().Add(-80 * time.Hour),
		time.Now().Add(-60 * time.Hour),
	}
	cfg := stale.DefaultConfig()
	results := stale.Detect(snaps, ts, cfg)
	if len(results) != 0 {
		t.Fatalf("expected no stale entries, got %d", len(results))
	}
}

func TestDetect_StaleEntry(t *testing.T) {
	e := makeEntry("DB_HOST", "old", "new", diff.StatusChanged)
	snaps := [][]diff.Entry{{e}, {e}, {e}, {e}}
	ts := []time.Time{
		time.Now().Add(-100 * time.Hour),
		time.Now().Add(-90 * time.Hour),
		time.Now().Add(-80 * time.Hour),
		time.Now().Add(-70 * time.Hour),
	}
	cfg := stale.DefaultConfig()
	results := stale.Detect(snaps, ts, cfg)
	if len(results) != 1 {
		t.Fatalf("expected 1 stale entry, got %d", len(results))
	}
	if results[0].Entry.Key != "DB_HOST" {
		t.Errorf("unexpected key: %s", results[0].Entry.Key)
	}
	if results[0].Repeats < 3 {
		t.Errorf("expected repeats >= 3, got %d", results[0].Repeats)
	}
}

func TestDetect_TooFewRepeats(t *testing.T) {
	e := makeEntry("X", "a", "b", diff.StatusChanged)
	snaps := [][]diff.Entry{{e}, {e}}
	ts := []time.Time{
		time.Now().Add(-100 * time.Hour),
		time.Now().Add(-90 * time.Hour),
	}
	cfg := stale.DefaultConfig() // MinRepeats = 3
	results := stale.Detect(snaps, ts, cfg)
	if len(results) != 0 {
		t.Fatalf("expected 0 stale entries, got %d", len(results))
	}
}

func TestDetect_TooRecent(t *testing.T) {
	e := makeEntry("Y", "a", "b", diff.StatusChanged)
	snaps := [][]diff.Entry{{e}, {e}, {e}}
	ts := []time.Time{
		time.Now().Add(-1 * time.Hour),
		time.Now().Add(-30 * time.Minute),
		time.Now().Add(-10 * time.Minute),
	}
	cfg := stale.DefaultConfig()
	results := stale.Detect(snaps, ts, cfg)
	if len(results) != 0 {
		t.Fatalf("expected 0 stale entries (too recent), got %d", len(results))
	}
}

func TestDetect_EmptySnapshots(t *testing.T) {
	results := stale.Detect(nil, nil, stale.DefaultConfig())
	if results != nil {
		t.Errorf("expected nil result for empty input")
	}
}

func TestDetect_CustomConfig(t *testing.T) {
	e := makeEntry("Z", "p", "q", diff.StatusAdded)
	snaps := [][]diff.Entry{{e}, {e}}
	ts := []time.Time{
		time.Now().Add(-5 * time.Hour),
		time.Now().Add(-3 * time.Hour),
	}
	cfg := stale.Config{MaxAge: 2 * time.Hour, MinRepeats: 2}
	results := stale.Detect(snaps, ts, cfg)
	if len(results) != 1 {
		t.Fatalf("expected 1 result with custom config, got %d", len(results))
	}
}
