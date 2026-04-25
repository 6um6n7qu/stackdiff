package retain

import (
	"testing"
	"time"

	"github.com/stackdiff/stackdiff/internal/diff"
)

func pipelineEntries(now time.Time) []diff.Entry {
	return []diff.Entry{
		makeEntry("changed-key", diff.StatusChanged, time.Hour, now),
		makeEntry("equal-key", diff.StatusEqual, time.Hour, now),
		makeEntry("added-key", diff.StatusAdded, 2*time.Hour, now),
	}
}

func TestRun_UsesDefaults(t *testing.T) {
	now := time.Now()
	entries := pipelineEntries(now)
	// Run uses DefaultConfig which excludes Equal
	kept := Run(entries)
	for _, e := range kept {
		if e.Status == diff.StatusEqual {
			t.Fatalf("equal entry should not be retained by default")
		}
	}
}

func TestRunWithConfig_CustomStatuses(t *testing.T) {
	now := time.Now()
	entries := pipelineEntries(now)
	cfg := Config{Statuses: []string{diff.StatusEqual}, Now: now}
	kept := RunWithConfig(entries, cfg)
	if len(kept) != 1 || kept[0].Key != "equal-key" {
		t.Fatalf("expected only equal-key, got %v", kept)
	}
}

func TestMustRetainSome_OK(t *testing.T) {
	now := time.Now()
	entries := pipelineEntries(now)
	cfg := DefaultConfig()
	cfg.Now = now
	cfg.MaxAge = 0
	if err := MustRetainSome(entries, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMustRetainSome_AllDropped(t *testing.T) {
	now := time.Now()
	entries := []diff.Entry{
		makeEntry("a", diff.StatusEqual, time.Hour, now),
	}
	cfg := DefaultConfig() // equal is not in default statuses
	cfg.Now = now
	cfg.MaxAge = 0
	if err := MustRetainSome(entries, cfg); err == nil {
		t.Fatal("expected error when all entries dropped")
	}
}

func TestMustRetainAll_NoneDropped(t *testing.T) {
	now := time.Now()
	entries := []diff.Entry{
		makeEntry("a", diff.StatusChanged, time.Hour, now),
	}
	cfg := DefaultConfig()
	cfg.Now = now
	cfg.MaxAge = 0
	if err := MustRetainAll(entries, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMustRetainAll_SomeDropped(t *testing.T) {
	now := time.Now()
	entries := pipelineEntries(now)
	cfg := DefaultConfig()
	cfg.Now = now
	cfg.MaxAge = 0
	// equal-key will be dropped
	if err := MustRetainAll(entries, cfg); err == nil {
		t.Fatal("expected error when some entries dropped")
	}
}
