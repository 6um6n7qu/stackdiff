package retain

import (
	"fmt"

	"github.com/stackdiff/stackdiff/internal/diff"
)

// Run applies the default retention policy to entries and returns kept entries.
func Run(entries []diff.Entry) []diff.Entry {
	return Apply(entries, DefaultConfig())
}

// RunWithConfig applies a custom retention policy.
func RunWithConfig(entries []diff.Entry, cfg Config) []diff.Entry {
	return Apply(entries, cfg)
}

// MustRetainSome returns an error if no entries survive the retention filter.
// This is useful as a pipeline guard to detect over-aggressive pruning.
func MustRetainSome(entries []diff.Entry, cfg Config) error {
	if len(entries) == 0 {
		return nil
	}
	kept := Apply(entries, cfg)
	if len(kept) == 0 {
		return fmt.Errorf("retain: all %d entries were dropped by retention policy", len(entries))
	}
	return nil
}

// MustRetainAll returns an error if any entry is dropped by the retention policy.
func MustRetainAll(entries []diff.Entry, cfg Config) error {
	kept := Apply(entries, cfg)
	if len(kept) != len(entries) {
		dropped := len(entries) - len(kept)
		return fmt.Errorf("retain: %d of %d entries dropped by retention policy", dropped, len(entries))
	}
	return nil
}
