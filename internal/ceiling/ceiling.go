// Package ceiling enforces an upper bound on the number of drift entries
// that flow through the pipeline. When the entry count exceeds the configured
// maximum, the list is trimmed and a sentinel entry is appended to signal
// truncation to downstream consumers.
package ceiling

import (
	"fmt"

	"github.com/your-org/stackdiff/internal/diff"
)

const defaultMax = 100

// Config controls the ceiling behaviour.
type Config struct {
	// Max is the maximum number of entries to allow through.
	// If zero, defaultMax is used.
	Max int

	// SentinelKey is the key used for the truncation notice entry.
	// Defaults to "__ceiling_truncated__".
	SentinelKey string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Max:         defaultMax,
		SentinelKey: "__ceiling_truncated__",
	}
}

// Apply trims entries to at most cfg.Max items. If truncation occurs a
// sentinel entry is appended so callers can detect the condition.
func Apply(entries []diff.Entry, cfg Config) []diff.Entry {
	if cfg.Max <= 0 {
		cfg.Max = defaultMax
	}
	if cfg.SentinelKey == "" {
		cfg.SentinelKey = "__ceiling_truncated__"
	}

	if len(entries) <= cfg.Max {
		return entries
	}

	truncated := make([]diff.Entry, cfg.Max, cfg.Max+1)
	copy(truncated, entries[:cfg.Max])

	sentinel := diff.Entry{
		Key:      cfg.SentinelKey,
		NewValue: fmt.Sprintf("truncated %d entries (max %d)", len(entries), cfg.Max),
		Status:   diff.StatusChanged,
	}
	truncated = append(truncated, sentinel)
	return truncated
}

// Exceeded reports whether the entry slice exceeds the configured maximum
// without modifying the slice.
func Exceeded(entries []diff.Entry, cfg Config) bool {
	if cfg.Max <= 0 {
		cfg.Max = defaultMax
	}
	return len(entries) > cfg.Max
}
