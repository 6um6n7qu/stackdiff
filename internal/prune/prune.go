// Package prune removes entries from a diff set based on age or staleness criteria.
package prune

import (
	"time"

	"github.com/stackdiff/internal/diff"
)

// Config controls pruning behaviour.
type Config struct {
	// MaxAge removes entries whose LastSeen is older than this duration.
	MaxAge time.Duration
	// Statuses limits pruning to entries with these statuses. Empty means all.
	Statuses []string
}

// DefaultConfig returns a Config that prunes entries older than 7 days.
func DefaultConfig() Config {
	return Config{
		MaxAge: 7 * 24 * time.Hour,
	}
}

// Result holds the outcome of a prune operation.
type Result struct {
	Kept    []diff.Entry
	Pruned  []diff.Entry
}

// Apply filters entries according to cfg, returning kept and pruned sets.
// Entries without metadata timestamp are always kept.
func Apply(entries []diff.Entry, cfg Config, now time.Time) Result {
	var result Result
	statusSet := toSet(cfg.Statuses)

	for _, e := range entries {
		if shouldPrune(e, cfg, statusSet, now) {
			result.Pruned = append(result.Pruned, e)
		} else {
			result.Kept = append(result.Kept, e)
		}
	}
	return result
}

func shouldPrune(e diff.Entry, cfg Config, statusSet map[string]bool, now time.Time) bool {
	if len(statusSet) > 0 && !statusSet[e.Status] {
		return false
	}
	raw, ok := e.Meta["last_seen"]
	if !ok {
		return false
	}
	ts, ok := raw.(time.Time)
	if !ok {
		return false
	}
	return now.Sub(ts) > cfg.MaxAge
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}
