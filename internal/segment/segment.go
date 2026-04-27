// Package segment splits a slice of diff entries into named buckets
// based on configurable predicate functions.
package segment

import (
	"fmt"

	"github.com/user/stackdiff/internal/diff"
)

// Predicate is a function that returns true if an entry belongs to a segment.
type Predicate func(e diff.Entry) bool

// Rule associates a named segment with a predicate.
type Rule struct {
	Name      string
	Predicate Predicate
}

// Result holds the segmented buckets and any entries that matched no rule.
type Result struct {
	Buckets   map[string][]diff.Entry
	Unmatched []diff.Entry
}

// HasDrift reports whether any bucket contains at least one drifting entry.
func (r Result) HasDrift() bool {
	for _, entries := range r.Buckets {
		for _, e := range entries {
			if e.IsDrift() {
				return true
			}
		}
	}
	return false
}

// BucketNames returns a sorted list of bucket names that contain entries.
func (r Result) BucketNames() []string {
	names := make([]string, 0, len(r.Buckets))
	for k := range r.Buckets {
		names = append(names, k)
	}
	return names
}

// Apply segments entries using the provided rules. Each entry is placed in the
// first matching bucket. Entries that match no rule are collected in Unmatched.
func Apply(entries []diff.Entry, rules []Rule) (Result, error) {
	if len(rules) == 0 {
		return Result{}, fmt.Errorf("segment: at least one rule is required")
	}

	buckets := make(map[string][]diff.Entry, len(rules))
	for _, r := range rules {
		buckets[r.Name] = nil
	}

	var unmatched []diff.Entry
	for _, e := range entries {
		placed := false
		for _, r := range rules {
			if r.Predicate(e) {
				buckets[r.Name] = append(buckets[r.Name], e)
				placed = true
				break
			}
		}
		if !placed {
			unmatched = append(unmatched, e)
		}
	}

	return Result{Buckets: buckets, Unmatched: unmatched}, nil
}

// DefaultRules returns a standard set of rules that segment by drift status.
func DefaultRules() []Rule {
	return []Rule{
		{Name: "changed", Predicate: func(e diff.Entry) bool { return e.Status == diff.StatusChanged }},
		{Name: "added", Predicate: func(e diff.Entry) bool { return e.Status == diff.StatusAdded }},
		{Name: "removed", Predicate: func(e diff.Entry) bool { return e.Status == diff.StatusRemoved }},
		{Name: "equal", Predicate: func(e diff.Entry) bool { return e.Status == diff.StatusEqual }},
	}
}
