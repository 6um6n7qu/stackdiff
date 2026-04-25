// Package expire provides TTL-based expiration filtering for diff entries.
// Entries whose metadata contains an "expires_at" timestamp in the past
// are considered expired and removed from the result set.
package expire

import (
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Now:        time.Now,
		MetaKey:    "expires_at",
		TimeFormat: time.RFC3339,
	}
}

// Config controls expiration behaviour.
type Config struct {
	// Now returns the current time; injectable for testing.
	Now func() time.Time
	// MetaKey is the metadata key that holds the expiry timestamp.
	MetaKey string
	// TimeFormat is the layout used to parse the timestamp value.
	TimeFormat string
}

// Apply removes entries whose expiry timestamp is in the past.
// Entries without the metadata key are kept unchanged.
func Apply(entries []diff.Entry, cfg Config) []diff.Entry {
	now := cfg.Now()
	out := make([]diff.Entry, 0, len(entries))
	for _, e := range entries {
		if isExpired(e, cfg, now) {
			continue
		}
		out = append(out, e)
	}
	return out
}

// CountExpired returns the number of entries that would be removed by Apply.
func CountExpired(entries []diff.Entry, cfg Config) int {
	now := cfg.Now()
	count := 0
	for _, e := range entries {
		if isExpired(e, cfg, now) {
			count++
		}
	}
	return count
}

func isExpired(e diff.Entry, cfg Config, now time.Time) bool {
	if e.Meta == nil {
		return false
	}
	raw, ok := e.Meta[cfg.MetaKey]
	if !ok || raw == "" {
		return false
	}
	t, err := time.Parse(cfg.TimeFormat, raw)
	if err != nil {
		return false
	}
	return t.Before(now)
}
