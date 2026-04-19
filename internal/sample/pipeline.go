package sample

import (
	"fmt"

	"github.com/user/stackdiff/internal/diff"
)

// Result holds sampled entries and metadata.
type Result struct {
	Entries  []diff.Entry
	Total    int
	Sampled  int
	Dropped  int
}

// String returns a human-readable summary of the sampling result.
func (r Result) String() string {
	return fmt.Sprintf("sampled %d/%d entries (%d dropped)", r.Sampled, r.Total, r.Dropped)
}

// Run applies sampling and returns a Result with metadata.
func Run(entries []diff.Entry, cfg Config) Result {
	sampled := Apply(entries, cfg)
	return Result{
		Entries: sampled,
		Total:   len(entries),
		Sampled: len(sampled),
		Dropped: len(entries) - len(sampled),
	}
}

// MustSample panics if sampling drops more than maxDropPct percent of entries.
func MustSample(entries []diff.Entry, cfg Config, maxDropPct float64) Result {
	r := Run(entries, cfg)
	if r.Total == 0 {
		return r
	}
	dropPct := float64(r.Dropped) / float64(r.Total) * 100
	if dropPct > maxDropPct {
		panic(fmt.Sprintf("sample: dropped %.1f%% exceeds allowed %.1f%%", dropPct, maxDropPct))
	}
	return r
}
