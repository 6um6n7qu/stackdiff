package threshold

import (
	"github.com/yourusername/stackdiff/internal/diff"
	"github.com/yourusername/stackdiff/internal/score"
)

// CheckEntries computes a drift score for the given entries and evaluates it
// against cfg, returning a Result.
func CheckEntries(entries []diff.Entry, cfg Config) Result {
	s := score.Compute(entries)
	return Evaluate(s, cfg)
}

// MustNotBreach returns an error string if the threshold is breached, or empty
// string otherwise. Useful for simple gate checks in pipelines.
func MustNotBreach(entries []diff.Entry, cfg Config) string {
	r := CheckEntries(entries, cfg)
	if r.Breached() {
		return r.Message
	}
	return ""
}
