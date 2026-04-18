// Package clamp provides utilities to enforce value length limits on config entries.
// Values exceeding the maximum length are truncated and marked.
package clamp

import (
	"fmt"

	"github.com/user/stackdiff/internal/diff"
)

const DefaultMaxLength = 256

// Config controls clamping behaviour.
type Config struct {
	MaxLength int
	Suffix    string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxLength: DefaultMaxLength,
		Suffix:    "...[truncated]",
	}
}

// Clamp returns a copy of entries where any value exceeding MaxLength is
// truncated and the configured suffix is appended.
func Clamp(entries []diff.Entry, cfg Config) []diff.Entry {
	if cfg.MaxLength <= 0 {
		cfg.MaxLength = DefaultMaxLength
	}
	out := make([]diff.Entry, len(entries))
	for i, e := range entries {
		e.OldValue = clampString(e.OldValue, cfg)
		e.NewValue = clampString(e.NewValue, cfg)
		out[i] = e
	}
	return out
}

func clampString(s string, cfg Config) string {
	if len(s) <= cfg.MaxLength {
		return s
	}
	return fmt.Sprintf("%s%s", s[:cfg.MaxLength], cfg.Suffix)
}
