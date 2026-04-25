// Package reorder provides utilities for sorting and reordering diff entries
// by various criteria such as key name, status priority, or custom comparators.
package reorder

import (
	"sort"
	"strings"

	"github.com/stackdiff/stackdiff/internal/diff"
)

// Strategy defines how entries should be reordered.
type Strategy string

const (
	// ByKey sorts entries alphabetically by key.
	ByKey Strategy = "key"
	// ByStatus sorts entries by drift status priority (changed > added > removed > equal).
	ByStatus Strategy = "status"
	// ByKeyDesc sorts entries reverse-alphabetically by key.
	ByKeyDesc Strategy = "key_desc"
)

// Config holds configuration for the reorder operation.
type Config struct {
	Strategy Strategy
	StableEqual bool // keep equal entries at the end
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Strategy:    ByStatus,
		StableEqual: true,
	}
}

var statusPriority = map[diff.Status]int{
	diff.StatusChanged: 0,
	diff.StatusAdded:   1,
	diff.StatusRemoved: 2,
	diff.StatusEqual:   3,
}

// Apply reorders entries according to the given Config.
// It does not mutate the input slice; a new slice is returned.
func Apply(entries []diff.Entry, cfg Config) []diff.Entry {
	if len(entries) == 0 {
		return entries
	}

	out := make([]diff.Entry, len(entries))
	copy(out, entries)

	switch cfg.Strategy {
	case ByKey:
		sort.SliceStable(out, func(i, j int) bool {
			return strings.ToLower(out[i].Key) < strings.ToLower(out[j].Key)
		})
	case ByKeyDesc:
		sort.SliceStable(out, func(i, j int) bool {
			return strings.ToLower(out[i].Key) > strings.ToLower(out[j].Key)
		})
	default: // ByStatus
		sort.SliceStable(out, func(i, j int) bool {
			pi := statusPriority[out[i].Status]
			pj := statusPriority[out[j].Status]
			if pi != pj {
				return pi < pj
			}
			return strings.ToLower(out[i].Key) < strings.ToLower(out[j].Key)
		})
	}

	return out
}
