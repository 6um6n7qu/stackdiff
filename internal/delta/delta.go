// Package delta computes a numeric delta between two snapshots,
// summarising how many keys were added, removed, or changed.
package delta

import (
	"fmt"

	"github.com/user/stackdiff/internal/diff"
)

// Result holds the numeric delta between two configs.
type Result struct {
	Added   int
	Removed int
	Changed int
	Total   int
}

// HasDrift returns true when any drift was detected.
func (r Result) HasDrift() bool {
	return r.Total > 0
}

// String returns a human-readable summary of the delta.
func (r Result) String() string {
	return fmt.Sprintf("added=%d removed=%d changed=%d total=%d", r.Added, r.Removed, r.Changed, r.Total)
}

// Compute derives a Result from a slice of diff entries.
func Compute(entries []diff.Entry) Result {
	var r Result
	for _, e := range entries {
		switch e.Status {
		case diff.StatusAdded:
			r.Added++
		case diff.StatusRemoved:
			r.Removed++
		case diff.StatusChanged:
			r.Changed++
		}
	}
	r.Total = r.Added + r.Removed + r.Changed
	return r
}
