// Package truncate provides utilities for truncating long config values
// in diff entries before display or export.
package truncate

import "github.com/user/stackdiff/internal/diff"

const DefaultMaxLen = 120

// Config holds truncation settings.
type Config struct {
	MaxLen int
	Suffix string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxLen: DefaultMaxLen,
		Suffix: "...",
	}
}

// Truncator applies value truncation to diff entries.
type Truncator struct {
	cfg Config
}

// New creates a Truncator with the given config.
func New(cfg Config) *Truncator {
	if cfg.MaxLen <= 0 {
		cfg.MaxLen = DefaultMaxLen
	}
	if cfg.Suffix == "" {
		cfg.Suffix = "..."
	}
	return &Truncator{cfg: cfg}
}

// Apply returns a new slice of entries with OldVal and NewVal truncated.
func (t *Truncator) Apply(entries []diff.Entry) []diff.Entry {
	out := make([]diff.Entry, len(entries))
	for i, e := range entries {
		e.OldVal = t.truncate(e.OldVal)
		e.NewVal = t.truncate(e.NewVal)
		out[i] = e
	}
	return out
}

func (t *Truncator) truncate(s string) string {
	if len(s) <= t.cfg.MaxLen {
		return s
	}
	return s[:t.cfg.MaxLen] + t.cfg.Suffix
}
