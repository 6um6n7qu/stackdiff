package ceiling_test

import (
	"strings"
	"testing"

	"github.com/your-org/stackdiff/internal/ceiling"
	"github.com/your-org/stackdiff/internal/diff"
)

func makeEntries(n int) []diff.Entry {
	out := make([]diff.Entry, n)
	for i := 0; i < n; i++ {
		out[i] = diff.Entry{
			Key:      strings.Repeat("k", i+1),
			NewValue: "v",
			Status:   diff.StatusChanged,
		}
	}
	return out
}

func TestApply_UnderLimit_ReturnsAll(t *testing.T) {
	entries := makeEntries(5)
	cfg := ceiling.DefaultConfig()
	cfg.Max = 10

	result := ceiling.Apply(entries, cfg)
	if len(result) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(result))
	}
}

func TestApply_ExactLimit_ReturnsAll(t *testing.T) {
	entries := makeEntries(10)
	cfg := ceiling.DefaultConfig()
	cfg.Max = 10

	result := ceiling.Apply(entries, cfg)
	if len(result) != 10 {
		t.Fatalf("expected 10 entries, got %d", len(result))
	}
}

func TestApply_OverLimit_TruncatesAndAppendsSentinel(t *testing.T) {
	entries := makeEntries(15)
	cfg := ceiling.DefaultConfig()
	cfg.Max = 10

	result := ceiling.Apply(entries, cfg)
	// max entries + 1 sentinel
	if len(result) != 11 {
		t.Fatalf("expected 11 entries, got %d", len(result))
	}

	sentinel := result[len(result)-1]
	if sentinel.Key != "__ceiling_truncated__" {
		t.Errorf("unexpected sentinel key: %q", sentinel.Key)
	}
	if !strings.Contains(sentinel.NewValue, "truncated") {
		t.Errorf("sentinel value should mention truncation, got %q", sentinel.NewValue)
	}
}

func TestApply_CustomSentinelKey(t *testing.T) {
	entries := makeEntries(5)
	cfg := ceiling.Config{Max: 3, SentinelKey: "__overflow__"}

	result := ceiling.Apply(entries, cfg)
	sentinel := result[len(result)-1]
	if sentinel.Key != "__overflow__" {
		t.Errorf("expected custom sentinel key, got %q", sentinel.Key)
	}
}

func TestApply_ZeroMax_UsesDefault(t *testing.T) {
	entries := makeEntries(5)
	cfg := ceiling.Config{} // Max == 0 → defaultMax (100)

	result := ceiling.Apply(entries, cfg)
	if len(result) != 5 {
		t.Fatalf("expected all 5 entries, got %d", len(result))
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries(5)
	cfg := ceiling.Config{Max: 3, SentinelKey: "__ceiling_truncated__"}

	_ = ceiling.Apply(entries, cfg)
	if len(entries) != 5 {
		t.Errorf("original slice was mutated")
	}
}

func TestExceeded_True(t *testing.T) {
	entries := makeEntries(11)
	cfg := ceiling.Config{Max: 10}
	if !ceiling.Exceeded(entries, cfg) {
		t.Error("expected Exceeded to return true")
	}
}

func TestExceeded_False(t *testing.T) {
	entries := makeEntries(10)
	cfg := ceiling.Config{Max: 10}
	if ceiling.Exceeded(entries, cfg) {
		t.Error("expected Exceeded to return false")
	}
}
