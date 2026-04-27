// Package split divides a slice of diff entries into two named partitions
// based on a user-supplied predicate. This is useful when a pipeline needs
// to route entries down separate processing paths (e.g. secrets vs. plain
// config) without losing any entries.
package split

import "github.com/stackdiff/stackdiff/internal/diff"

// Result holds the two partitions produced by Apply.
type Result struct {
	// Name is the label given to this partition set.
	Name string
	// Matched contains entries for which the predicate returned true.
	Matched []diff.Entry
	// Unmatched contains entries for which the predicate returned false.
	Unmatched []diff.Entry
}

// HasDrift reports whether either partition contains a drifting entry.
func (r Result) HasDrift() bool {
	for _, e := range r.Matched {
		if e.IsDrift() {
			return true
		}
	}
	for _, e := range r.Unmatched {
		if e.IsDrift() {
			return true
		}
	}
	return false
}

// MatchedDrift returns only drifting entries from the Matched partition.
func (r Result) MatchedDrift() []diff.Entry {
	out := make([]diff.Entry, 0, len(r.Matched))
	for _, e := range r.Matched {
		if e.IsDrift() {
			out = append(out, e)
		}
	}
	return out
}

// Apply partitions entries into Matched and Unmatched using pred.
// All entries are preserved; none are dropped.
func Apply(name string, entries []diff.Entry, pred func(diff.Entry) bool) Result {
	r := Result{Name: name}
	for _, e := range entries {
		if pred(e) {
			r.Matched = append(r.Matched, e)
		} else {
			r.Unmatched = append(r.Unmatched, e)
		}
	}
	return r
}

// ByStatus is a convenience predicate factory that matches entries whose
// Status equals one of the supplied statuses.
func ByStatus(statuses ...string) func(diff.Entry) bool {
	set := make(map[string]struct{}, len(statuses))
	for _, s := range statuses {
		set[s] = struct{}{}
	}
	return func(e diff.Entry) bool {
		_, ok := set[e.Status]
		return ok
	}
}
