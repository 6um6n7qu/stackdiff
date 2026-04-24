package quota

import (
	"fmt"

	"github.com/user/stackdiff/internal/diff"
)

// Run applies quota enforcement to entries using the default configuration.
// Entries that exceed the per-key limit within the window are dropped.
func Run(entries []diff.Entry) []diff.Entry {
	return New(DefaultConfig()).Filter(entries)
}

// RunWithConfig applies quota enforcement using the provided Config.
func RunWithConfig(cfg Config, entries []diff.Entry) []diff.Entry {
	return New(cfg).Filter(entries)
}

// MustAllow returns an error if any entry in the slice would be blocked by the
// enforcer. It is useful in strict pipeline stages where quota breach is fatal.
func MustAllow(e *Enforcer, entries []diff.Entry) error {
	for _, en := range entries {
		if !e.Allow(en) {
			return fmt.Errorf("quota exceeded for key %q", en.Key)
		}
	}
	return nil
}
