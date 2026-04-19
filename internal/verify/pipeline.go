package verify

import (
	"fmt"
	"io"

	"github.com/stackdiff/stackdiff/internal/diff"
)

// Run evaluates entries and writes a human-readable summary to w.
// Returns an error if verification fails.
func Run(entries []diff.Entry, cfg Config, w io.Writer) error {
	result := Evaluate(entries, cfg)

	for _, msg := range result.Messages {
		fmt.Fprintln(w, msg)
	}

	switch result.Status {
	case StatusPass:
		fmt.Fprintln(w, "verification passed: no drift detected")
		return nil
	case StatusWarning:
		fmt.Fprintf(w, "verification warning: %d drift entries (within tolerance)\n", len(result.Entries))
		return nil
	default:
		return fmt.Errorf("verification failed: %d drift entries detected", len(result.Entries))
	}
}

// MustPass is like Run but panics on failure (useful in test helpers).
func MustPass(entries []diff.Entry, cfg Config, w io.Writer) {
	if err := Run(entries, cfg, w); err != nil {
		panic(err)
	}
}
