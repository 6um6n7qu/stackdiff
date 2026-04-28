package watermark

import (
	"fmt"

	"github.com/user/stackdiff/internal/diff"
)

// Result holds the outcome of a pipeline check.
type Result struct {
	Series   string
	Count    int
	Peak     int
	NewPeak  bool
}

// String returns a human-readable summary.
func (r Result) String() string {
	if r.NewPeak {
		return fmt.Sprintf("series %q: new peak %d (previous %d)", r.Series, r.Count, r.Peak)
	}
	return fmt.Sprintf("series %q: count %d within peak %d", r.Series, r.Count, r.Peak)
}

// Run records the drift count from entries against the store and returns a Result.
func Run(s *Store, series string, entries []diff.Entry) Result {
	count := 0
	for _, e := range entries {
		if e.IsDrift() {
			count++
		}
	}

	prev, _ := s.Get(series)
	newPeak := s.Record(series, count)

	return Result{
		Series:  series,
		Count:   count,
		Peak:    prev.Peak,
		NewPeak: newPeak,
	}
}

// MustNotExceedPeak returns an error if the current drift count exceeds the
// stored peak by more than tolerance.
func MustNotExceedPeak(s *Store, series string, entries []diff.Entry, tolerance int) error {
	r := Run(s, series, entries)
	if r.NewPeak && r.Count > r.Peak+tolerance {
		return fmt.Errorf("watermark: drift count %d exceeds previous peak %d by more than tolerance %d",
			r.Count, r.Peak, tolerance)
	}
	return nil
}
