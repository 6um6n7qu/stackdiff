package quota_test

import (
	"testing"
	"time"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/quota"
)

func pipelineEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "A", NewValue: "1", Status: diff.StatusChanged},
		{Key: "A", NewValue: "2", Status: diff.StatusChanged},
		{Key: "A", NewValue: "3", Status: diff.StatusChanged},
		{Key: "B", NewValue: "x", Status: diff.StatusAdded},
	}
}

func TestRun_DefaultDropsOverQuota(t *testing.T) {
	// default MaxPerKey=5, so all 4 entries should pass with default config
	got := quota.Run(pipelineEntries())
	if len(got) != 4 {
		t.Fatalf("expected 4, got %d", len(got))
	}
}

func TestRunWithConfig_DropsExcess(t *testing.T) {
	cfg := quota.Config{MaxPerKey: 2, Window: time.Minute}
	got := quota.RunWithConfig(cfg, pipelineEntries())
	// A appears 3 times but limit is 2; B appears 1 time
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}

func TestMustAllow_ReturnsNilWhenUnderQuota(t *testing.T) {
	e := quota.New(quota.Config{MaxPerKey: 10, Window: time.Minute})
	entries := []diff.Entry{
		{Key: "X", Status: diff.StatusChanged},
	}
	if err := quota.MustAllow(e, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMustAllow_ReturnsErrorWhenExceeded(t *testing.T) {
	e := quota.New(quota.Config{MaxPerKey: 1, Window: time.Minute})
	entry := diff.Entry{Key: "Y", Status: diff.StatusChanged}
	e.Allow(entry) // consume the one allowed slot
	entries := []diff.Entry{entry}
	if err := quota.MustAllow(e, entries); err == nil {
		t.Fatal("expected error for exceeded quota")
	}
}
