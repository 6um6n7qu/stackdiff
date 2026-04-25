package retain

import (
	"testing"
	"time"

	"github.com/stackdiff/stackdiff/internal/diff"
)

func makeEntry(key, status string, age time.Duration, now time.Time) diff.Entry {
	e := diff.Entry{
		Key:    key,
		Status: status,
		Meta:   map[string]string{},
	}
	if age >= 0 {
		e.Meta["timestamp"] = now.Add(-age).Format(time.RFC3339)
	}
	return e
}

func TestApply_KeepsMatchingStatus(t *testing.T) {
	now := time.Now()
	entries := []diff.Entry{
		makeEntry("a", diff.StatusChanged, time.Hour, now),
		makeEntry("b", diff.StatusEqual, time.Hour, now),
	}
	cfg := DefaultConfig()
	cfg.Now = now
	cfg.MaxAge = 0
	kept := Apply(entries, cfg)
	if len(kept) != 1 || kept[0].Key != "a" {
		t.Fatalf("expected only 'a', got %v", kept)
	}
}

func TestApply_DropsOldEntries(t *testing.T) {
	now := time.Now()
	entries := []diff.Entry{
		makeEntry("recent", diff.StatusChanged, 1*time.Hour, now),
		makeEntry("old", diff.StatusChanged, 100*time.Hour, now),
	}
	cfg := DefaultConfig()
	cfg.Now = now
	cfg.MaxAge = 24 * time.Hour
	kept := Apply(entries, cfg)
	if len(kept) != 1 || kept[0].Key != "recent" {
		t.Fatalf("expected only 'recent', got %v", kept)
	}
}

func TestApply_NoTimestamp_AlwaysKept(t *testing.T) {
	now := time.Now()
	e := diff.Entry{Key: "x", Status: diff.StatusAdded, Meta: map[string]string{}}
	cfg := DefaultConfig()
	cfg.Now = now
	cfg.MaxAge = time.Minute
	kept := Apply([]diff.Entry{e}, cfg)
	if len(kept) != 1 {
		t.Fatal("entry without timestamp should always be kept")
	}
}

func TestApply_EmptyStatuses_KeepsAll(t *testing.T) {
	now := time.Now()
	entries := []diff.Entry{
		makeEntry("a", diff.StatusChanged, time.Hour, now),
		makeEntry("b", diff.StatusEqual, time.Hour, now),
		makeEntry("c", diff.StatusAdded, time.Hour, now),
	}
	cfg := Config{Statuses: nil, MaxAge: 0, Now: now}
	kept := Apply(entries, cfg)
	if len(kept) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(kept))
	}
}

func TestCountRetained(t *testing.T) {
	now := time.Now()
	entries := []diff.Entry{
		makeEntry("a", diff.StatusChanged, time.Hour, now),
		makeEntry("b", diff.StatusEqual, time.Hour, now),
	}
	cfg := DefaultConfig()
	cfg.Now = now
	cfg.MaxAge = 0
	if got := CountRetained(entries, cfg); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}
