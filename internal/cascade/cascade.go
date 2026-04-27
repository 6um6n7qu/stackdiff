// Package cascade propagates diff entries through a chain of named stages,
// collecting per-stage results for pipeline observability.
package cascade

import (
	"fmt"

	"github.com/yourusername/stackdiff/internal/diff"
)

// Stage is a named transformation applied to a slice of diff entries.
type Stage struct {
	Name string
	Fn   func([]diff.Entry) ([]diff.Entry, error)
}

// Result holds the output of a single stage.
type Result struct {
	Stage   string
	Entries []diff.Entry
	Dropped int
}

// Report is the full output of a cascade run.
type Report struct {
	Stages []Result
	Final  []diff.Entry
}

// HasDrift reports whether the final entry set contains any drift.
func (r *Report) HasDrift() bool {
	for _, e := range r.Final {
		if e.IsDrift() {
			return true
		}
	}
	return false
}

// Run executes each stage in order, feeding the output of one into the next.
// If any stage returns an error the run is aborted and the error is returned.
func Run(entries []diff.Entry, stages []Stage) (*Report, error) {
	if len(stages) == 0 {
		return &Report{Final: entries}, nil
	}

	report := &Report{}
	current := entries

	for _, s := range stages {
		if s.Fn == nil {
			return nil, fmt.Errorf("cascade: stage %q has nil function", s.Name)
		}
		next, err := s.Fn(current)
		if err != nil {
			return nil, fmt.Errorf("cascade: stage %q: %w", s.Name, err)
		}
		report.Stages = append(report.Stages, Result{
			Stage:   s.Name,
			Entries: next,
			Dropped: len(current) - len(next),
		})
		current = next
	}

	report.Final = current
	return report, nil
}
