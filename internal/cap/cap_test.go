package cap_test

import (
	"strings"
	"testing"

	"github.com/user/stackdiff/internal/cap"
	"github.com/user/stackdiff/internal/diff"
)

func makeEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "A", NewValue: "1", Status: diff.StatusChanged},
		{Key: "B", NewValue: "2", Status: diff.StatusChanged},
		{Key: "C", NewValue: "3", Status: diff.StatusChanged},
		{Key: "D", NewValue: "4", Status: diff.StatusAdded},
		{Key: "E", NewValue: "5", Status: diff.StatusAdded},
		{Key: "F", OldValue: "6", Status: diff.StatusRemoved},
	}
}

func TestApply_UnderLimit_ReturnsAll(t *testing.T) {
	entries := makeEntries()
	cfg := cap.Config{MaxPerStatus: 5}
	out := cap.Apply(entries, cfg)
	// No truncation expected; sentinel count should be zero.
	for _, e := range out {
		if e.Key == cfg.SentinelKey {
			t.Errorf("unexpected sentinel entry")
		}
	}
	if len(out) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(out))
	}
}

func TestApply_TruncatesChangedBucket(t *testing.T) {
	entries := makeEntries() // 3 changed
	cfg := cap.DefaultConfig()
	cfg.MaxPerStatus = 2
	out := cap.Apply(entries, cfg)

	changedCount := 0
	for _, e := range out {
		if e.Status == diff.StatusChanged && e.Key != cfg.SentinelKey {
			changedCount++
		}
	}
	if changedCount != 2 {
		t.Errorf("expected 2 changed entries, got %d", changedCount)
	}
}

func TestApply_AppendsSentinel(t *testing.T) {
	entries := makeEntries() // 3 changed
	cfg := cap.DefaultConfig()
	cfg.MaxPerStatus = 1
	out := cap.Apply(entries, cfg)

	var sentinel *diff.Entry
	for i := range out {
		if out[i].Key == cfg.SentinelKey && out[i].Status == diff.StatusChanged {
			sentinel = &out[i]
			break
		}
	}
	if sentinel == nil {
		t.Fatal("expected sentinel entry for changed bucket")
	}
	if !strings.Contains(sentinel.NewValue, "truncated") {
		t.Errorf("sentinel value should mention truncated, got %q", sentinel.NewValue)
	}
}

func TestApply_ZeroMax_NoLimit(t *testing.T) {
	entries := makeEntries()
	cfg := cap.Config{MaxPerStatus: 0}
	out := cap.Apply(entries, cfg)
	if len(out) != len(entries) {
		t.Errorf("expected all entries with zero max, got %d", len(out))
	}
}

func TestTruncated_True(t *testing.T) {
	entries := makeEntries()
	cfg := cap.Config{MaxPerStatus: 1}
	if !cap.Truncated(entries, cfg) {
		t.Error("expected Truncated to return true")
	}
}

func TestTruncated_False(t *testing.T) {
	entries := makeEntries()
	cfg := cap.Config{MaxPerStatus: 10}
	if cap.Truncated(entries, cfg) {
		t.Error("expected Truncated to return false")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries()
	orig := make([]diff.Entry, len(entries))
	copy(orig, entries)
	cfg := cap.Config{MaxPerStatus: 1}
	cap.Apply(entries, cfg)
	for i, e := range entries {
		if e != orig[i] {
			t.Errorf("original slice mutated at index %d", i)
		}
	}
}
