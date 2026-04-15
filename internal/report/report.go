package report

import (
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// Format represents the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatMarkdown Format = "markdown"
)

// Report holds metadata and diff results for a comparison run.
type Report struct {
	GeneratedAt time.Time        `json:"generated_at"`
	SourceLabel string           `json:"source_label"`
	TargetLabel string           `json:"target_label"`
	Diffs       []diff.DiffEntry `json:"diffs"`
	Summary     Summary          `json:"summary"`
}

// Summary contains aggregated counts from the diff.
type Summary struct {
	Added   int `json:"added"`
	Removed int `json:"removed"`
	Changed int `json:"changed"`
	Total   int `json:"total"`
}

// New creates a new Report from a slice of DiffEntry values.
func New(sourceLabel, targetLabel string, diffs []diff.DiffEntry) *Report {
	s := Summary{Total: len(diffs)}
	for _, d := range diffs {
		switch d.Kind {
		case diff.Added:
			s.Added++
		case diff.Removed:
			s.Removed++
		case diff.Changed:
			s.Changed++
		}
	}
	return &Report{
		GeneratedAt: time.Now().UTC(),
		SourceLabel: sourceLabel,
		TargetLabel: targetLabel,
		Diffs:       diffs,
		Summary:     s,
	}
}

// HasDrift returns true when the report contains at least one diff entry.
func (r *Report) HasDrift() bool {
	return r.Summary.Total > 0
}
