// Package coalesce provides utilities for merging multiple sets of diff entries
// by selecting the first non-empty (or non-equal) value for each key across
// an ordered list of entry slices. This is useful when layering config sources
// such as defaults, environment overrides, and live values.
package coalesce

import (
	"github.com/yourusername/stackdiff/internal/diff"
)

// Strategy controls how coalescing picks a winner when multiple sources
// provide a value for the same key.
type Strategy int

const (
	// StrategyFirstNonEmpty selects the first source whose NewValue is non-empty.
	StrategyFirstNonEmpty Strategy = iota

	// StrategyFirstDrift selects the first source whose entry is not StatusEqual.
	StrategyFirstDrift

	// StrategyLast always uses the last source that contains the key.
	StrategyLast
)

// Config holds options for the coalesce operation.
type Config struct {
	// Strategy determines how a winner is chosen per key.
	Strategy Strategy

	// FallbackToEqual, when true, includes keys whose only available value
	// is StatusEqual (i.e. no drift found in any source).
	FallbackToEqual bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Strategy:        StrategyFirstNonEmpty,
		FallbackToEqual: true,
	}
}

// Apply coalesces multiple slices of diff.Entry into a single slice.
// Sources are evaluated in order; the strategy determines which source
// "wins" when a key appears in more than one slice.
func Apply(sources [][]diff.Entry, cfg Config) []diff.Entry {
	if len(sources) == 0 {
		return nil
	}

	// Track insertion order so output is deterministic.
	seen := make(map[string]struct{})
	order := []string{}
	result := make(map[string]diff.Entry)

	for _, entries := range sources {
		for _, e := range entries {
			if _, exists := seen[e.Key]; !exists {
				seen[e.Key] = struct{}{}
				order = append(order, e.Key)
			}
		}
	}

	for _, key := range order {
		winner, ok := pick(key, sources, cfg)
		if ok {
			result[key] = winner
		}
	}

	out := make([]diff.Entry, 0, len(order))
	for _, key := range order {
		if e, ok := result[key]; ok {
			out = append(out, e)
		}
	}
	return out
}

// Count returns the number of entries that would be produced by Apply.
func Count(sources [][]diff.Entry, cfg Config) int {
	return len(Apply(sources, cfg))
}

// pick selects the winning entry for key across all sources according to cfg.
func pick(key string, sources [][]diff.Entry, cfg Config) (diff.Entry, bool) {
	var candidates []diff.Entry
	for _, entries := range sources {
		for _, e := range entries {
			if e.Key == key {
				candidates = append(candidates, e)
				break
			}
		}
	}

	if len(candidates) == 0 {
		return diff.Entry{}, false
	}

	switch cfg.Strategy {
	case StrategyFirstDrift:
		for _, c := range candidates {
			if c.Status != diff.StatusEqual {
				return c, true
			}
		}
		if cfg.FallbackToEqual {
			return candidates[0], true
		}
		return diff.Entry{}, false

	case StrategyLast:
		return candidates[len(candidates)-1], true

	default: // StrategyFirstNonEmpty
		for _, c := range candidates {
			if c.NewValue != "" {
				return c, true
			}
		}
		if cfg.FallbackToEqual {
			return candidates[0], true
		}
		return diff.Entry{}, false
	}
}
