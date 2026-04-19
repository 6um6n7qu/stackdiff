package diff

import (
	"testing"
)

func makeTestEntry(key, oldVal, newVal, status string) Entry {
	return Entry{Key: key, OldVal: oldVal, NewVal: newVal, Status: status}
}

func TestMergeEntries_RightStrategyOverwrites(t *testing.T) {
	base := []Entry{makeTestEntry("PORT", "8080", "8080", StatusEqual)}
	overlay := []Entry{makeTestEntry("PORT", "8080", "9090", StatusChanged)}

	result := MergeEntries(base, overlay, DefaultMergeOptions())
	if len(result) != 1 || result[0].NewVal != "9090" {
		t.Errorf("expected overlay value 9090, got %v", result)
	}
}

func TestMergeEntries_LeftStrategyKeepsBase(t *testing.T) {
	base := []Entry{makeTestEntry("PORT", "8080", "8080", StatusEqual)}
	overlay := []Entry{makeTestEntry("PORT", "8080", "9090", StatusChanged)}

	result := MergeEntries(base, overlay, MergeOptions{Strategy: MergeStrategyLeft})
	if len(result) != 1 || result[0].NewVal != "8080" {
		t.Errorf("expected base value 8080, got %v", result)
	}
}

func TestMergeEntries_UniqueKeysIncluded(t *testing.T) {
	base := []Entry{makeTestEntry("HOST", "localhost", "localhost", StatusEqual)}
	overlay := []Entry{makeTestEntry("PORT", "", "8080", StatusAdded)}

	result := MergeEntries(base, overlay, DefaultMergeOptions())
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestMergeEntries_EmptyOverlay(t *testing.T) {
	base := []Entry{makeTestEntry("HOST", "localhost", "localhost", StatusEqual)}
	result := MergeEntries(base, nil, DefaultMergeOptions())
	if len(result) != 1 || result[0].Key != "HOST" {
		t.Errorf("expected base entry, got %v", result)
	}
}

func TestMergeEntries_EmptyBase(t *testing.T) {
	overlay := []Entry{makeTestEntry("PORT", "", "8080", StatusAdded)}
	result := MergeEntries(nil, overlay, DefaultMergeOptions())
	if len(result) != 1 || result[0].Key != "PORT" {
		t.Errorf("expected overlay entry, got %v", result)
	}
}

func TestDefaultMergeOptions_IsRight(t *testing.T) {
	opts := DefaultMergeOptions()
	if opts.Strategy != MergeStrategyRight {
		t.Errorf("expected MergeStrategyRight, got %v", opts.Strategy)
	}
}
