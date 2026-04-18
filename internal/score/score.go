// Package score computes a numeric drift severity score from diff entries.
package score

import (
	"fmt"

	"github.com/yourorg/stackdiff/internal/diff"
)

// Weights assigned to each drift status.
const (
	WeightChanged = 2
	WeightAdded   = 1
	WeightRemoved = 1
)

// Result holds the computed score and a breakdown by status.
type Result struct {
	Total   int
	Changed int
	Added   int
	Removed int
}

// Grade returns a letter grade based on the total score.
func (r Result) Grade() string {
	switch {
	case r.Total == 0:
		return "A"
	case r.Total <= 3:
		return "B"
	case r.Total <= 8:
		return "C"
	case r.Total <= 15:
		return "D"
	default:
		return "F"
	}
}

// String returns a human-readable summary.
func (r Result) String() string {
	return fmt.Sprintf("score=%d grade=%s (changed=%d added=%d removed=%d)",
		r.Total, r.Grade(), r.Changed, r.Added, r.Removed)
}

// Compute calculates a drift severity score from the given entries.
func Compute(entries []diff.Entry) Result {
	var r Result
	for _, e := range entries {
		switch e.Status {
		case diff.StatusChanged:
			r.Changed++
			r.Total += WeightChanged
		case diff.StatusAdded:
			r.Added++
			r.Total += WeightAdded
		case diff.StatusRemoved:
			r.Removed++
			r.Total += WeightRemoved
		}
	}
	return r
}
