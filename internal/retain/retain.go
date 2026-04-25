// Package retain provides filtering logic to keep only entries that
// match a specified retention policy based on status and age.
package retain

import (
	"time"

	"github.com/stackdiff/stackdiff/internal/diff"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxAge:   72 * time.Hour,
		Statuses: []string{diff.StatusChanged, diff.StatusAdded, diff.StatusRemoved},
	}
}

// Config controls which entries are retained.
type Config struct {
	// MaxAge is the maximum age of an entry to retain. Zero means no limit.
	MaxAge time.Duration
	// Statuses lists the drift statuses to retain. Empty means retain all.
	Statuses []string
	// Now overrides the current time (useful for testing).
	Now time.Time
}

// Apply filters entries according to the retention policy and returns
// only those that should be kept.
func Apply(entries []diff.Entry, cfg Config) []diff.Entry {
	now := cfg.Now
	if now.IsZero() {
		now = time.Now()
	}

	statusSet := make(map[string]struct{}, len(cfg.Statuses))
	for _, s := range cfg.Statuses {
		statusSet[s] = struct{}{}
	}

	var kept []diff.Entry
	for _, e := range entries {
		if len(statusSet) > 0 {
			if _, ok := statusSet[e.Status]; !ok {
				continue
			}
		}
		if cfg.MaxAge > 0 {
			ts, ok := e.Meta["timestamp"]
			if ok {
				t, err := time.Parse(time.RFC3339, ts)
				if err == nil && now.Sub(t) > cfg.MaxAge {
					continue
				}
			}
		}
		kept = append(kept, e)
	}
	return kept
}

// CountRetained returns the number of entries that would be retained
// under the given config without allocating a new slice.
func CountRetained(entries []diff.Entry, cfg Config) int {
	return len(Apply(entries, cfg))
}
