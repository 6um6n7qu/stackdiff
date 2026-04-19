// Package sample provides utilities for sampling drift entries
// based on rate or count limits, useful for high-volume comparisons.
package sample

import (
	"math/rand"

	"github.com/user/stackdiff/internal/diff"
)

// Config controls how sampling is applied.
type Config struct {
	// Rate is a value between 0.0 and 1.0 representing the fraction to keep.
	// A rate of 1.0 keeps all entries.
	Rate float64
	// MaxEntries caps the total number of entries returned. 0 means no cap.
	MaxEntries int
	// Seed is used for reproducible sampling. 0 means random.
	Seed int64
}

// DefaultConfig returns a Config that keeps all entries.
func DefaultConfig() Config {
	return Config{Rate: 1.0}
}

// Apply samples entries according to cfg.
func Apply(entries []diff.Entry, cfg Config) []diff.Entry {
	rate := cfg.Rate
	if rate <= 0 {
		return []diff.Entry{}
	}
	if rate >= 1.0 && cfg.MaxEntries == 0 {
		return entries
	}

	var rng *rand.Rand
	if cfg.Seed != 0 {
		rng = rand.New(rand.NewSource(cfg.Seed))
	} else {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}

	out := make([]diff.Entry, 0, len(entries))
	for _, e := range entries {
		if rate >= 1.0 || rng.Float64() < rate {
			out = append(out, e)
		}
	}

	if cfg.MaxEntries > 0 && len(out) > cfg.MaxEntries {
		out = out[:cfg.MaxEntries]
	}
	return out
}

// Count returns how many entries would be sampled without allocating a full slice.
func Count(entries []diff.Entry, cfg Config) int {
	return len(Apply(entries, cfg))
}
