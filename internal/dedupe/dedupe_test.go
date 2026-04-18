package dedupe_test

import (
	"testing"

	"github.com/yourorg/stackdiff/internal/dedupe"
	"github.com/yourorg/stackdiff/internal/diff"
)

func makeEntry(key, old, new string, status diff.Status) diff.Entry {
	return diff.Entry{Key: key, OldValue: old, NewValue: new, Status: status}
}

func TestApply_NoDuplicates(t *testing.T) {
	entries := []diff.Entry{
		makeEntry("A", "", "1", diff.StatusAdded),
		makeEntry("B", "x", "y", diff.StatusChanged),
	}
	result := dedupe.Apply(entries, dedupe.StrategyFirst)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestApply_StrategyFirst(t *testing.T) {
	entries := []diff.Entry{
		makeEntry("KEY", "", "first", diff.StatusAdded),
		makeEntry("KEY", "", "second", diff.StatusAdded),
	}
	result := dedupe.Apply(entries, dedupe.StrategyFirst)
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].NewValue != "first" {
		t.Errorf("expected 'first', got %q", result[0].NewValue)
	}
}

func TestApply_StrategyLast(t *testing.T) {
	entries := []diff.Entry{
		makeEntry("KEY", "", "first", diff.StatusAdded),
		makeEntry("KEY", "", "second", diff.StatusAdded),
	}
	result := dedupe.Apply(entries, dedupe.StrategyLast)
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].NewValue != "second" {
		t.Errorf("expected 'second', got %q", result[0].NewValue)
	}
}

func TestApply_PreservesOrder(t *testing.T) {
	entries := []diff.Entry{
		makeEntry("A", "", "1", diff.StatusAdded),
		makeEntry("B", "", "2", diff.StatusAdded),
		makeEntry("A", "", "3", diff.StatusAdded),
	}
	result := dedupe.Apply(entries, dedupe.StrategyFirst)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0].Key != "A" || result[1].Key != "B" {
		t.Errorf("order not preserved: got %v %v", result[0].Key, result[1].Key)
	}
}

func TestApply_Empty(t *testing.T) {
	result := dedupe.Apply(nil, dedupe.StrategyFirst)
	if len(result) != 0 {
		t.Errorf("expected empty result")
	}
}

func TestCount_NoDuplicates(t *testing.T) {
	entries := []diff.Entry{
		makeEntry("A", "", "1", diff.StatusAdded),
		makeEntry("B", "", "2", diff.StatusAdded),
	}
	if n := dedupe.Count(entries); n != 0 {
		t.Errorf("expected 0 dupes, got %d", n)
	}
}

func TestCount_WithDuplicates(t *testing.T) {
	entries := []diff.Entry{
		makeEntry("A", "", "1", diff.StatusAdded),
		makeEntry("A", "", "2", diff.StatusAdded),
		makeEntry("A", "", "3", diff.StatusAdded),
	}
	if n := dedupe.Count(entries); n != 2 {
		t.Errorf("expected 2 dupes, got %d", n)
	}
}
