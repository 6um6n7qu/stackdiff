package sample_test

import (
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/sample"
)

func makeEntries(n int) []diff.Entry {
	out := make([]diff.Entry, n)
	for i := 0; i < n; i++ {
		out[i] = diff.Entry{
			Key:      "KEY",
			OldValue: "old",
			NewValue: "new",
			Status:   diff.StatusChanged,
		}
	}
	return out
}

func TestApply_RateOne_ReturnsAll(t *testing.T) {
	entries := makeEntries(20)
	result := sample.Apply(entries, sample.DefaultConfig())
	if len(result) != 20 {
		t.Fatalf("expected 20, got %d", len(result))
	}
}

func TestApply_RateZero_ReturnsEmpty(t *testing.T) {
	entries := makeEntries(20)
	result := sample.Apply(entries, sample.Config{Rate: 0})
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}

func TestApply_MaxEntries_Caps(t *testing.T) {
	entries := makeEntries(50)
	cfg := sample.Config{Rate: 1.0, MaxEntries: 10}
	result := sample.Apply(entries, cfg)
	if len(result) != 10 {
		t.Fatalf("expected 10, got %d", len(result))
	}
}

func TestApply_Reproducible_WithSeed(t *testing.T) {
	entries := makeEntries(100)
	cfg := sample.Config{Rate: 0.5, Seed: 99}
	r1 := sample.Apply(entries, cfg)
	r2 := sample.Apply(entries, cfg)
	if len(r1) != len(r2) {
		t.Fatalf("expected same count with same seed: %d vs %d", len(r1), len(r2))
	}
}

func TestApply_PartialRate_ReducesCount(t *testing.T) {
	entries := makeEntries(1000)
	cfg := sample.Config{Rate: 0.1, Seed: 1}
	result := sample.Apply(entries, cfg)
	if len(result) > 200 {
		t.Fatalf("expected roughly 100, got %d", len(result))
	}
}

func TestCount_MatchesApply(t *testing.T) {
	entries := makeEntries(50)
	cfg := sample.Config{Rate: 1.0, MaxEntries: 20}
	if sample.Count(entries, cfg) != len(sample.Apply(entries, cfg)) {
		t.Fatal("Count and Apply disagree")
	}
}

func TestDefaultConfig_RateIsOne(t *testing.T) {
	cfg := sample.DefaultConfig()
	if cfg.Rate != 1.0 {
		t.Fatalf("expected rate 1.0, got %f", cfg.Rate)
	}
}
