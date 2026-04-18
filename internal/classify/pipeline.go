package classify

import "github.com/yourusername/stackdiff/internal/diff"

// BySeverity filters Results keeping only those at or above the given severity.
// Severity order: critical > high > medium > low.
func BySeverity(results []Result, min Severity) []Result {
	rank := map[Severity]int{
		SeverityLow:      0,
		SeverityMedium:   1,
		SeverityHigh:     2,
		SeverityCritical: 3,
	}
	minRank := rank[min]
	out := make([]Result, 0)
	for _, r := range results {
		if rank[r.Severity] >= minRank {
			out = append(out, r)
		}
	}
	return out
}

// Entries extracts the diff.Entry slice from a []Result.
func Entries(results []Result) []diff.Entry {
	out := make([]diff.Entry, len(results))
	for i, r := range results {
		out[i] = r.Entry
	}
	return out
}

// ClassifyAndFilter is a convenience function that classifies entries and
// returns only those at or above the minimum severity.
func ClassifyAndFilter(entries []diff.Entry, rules []Rule, min Severity) []Result {
	c := New(rules)
	return BySeverity(c.Apply(entries), min)
}
