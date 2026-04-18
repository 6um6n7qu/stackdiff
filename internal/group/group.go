// Package group provides utilities for grouping diff entries by a shared attribute.
package group

import (
	"sort"

	"github.com/yourorg/stackdiff/internal/diff"
)

// GroupFunc extracts a group key from an entry.
type GroupFunc func(e diff.Entry) string

// Group holds a named collection of entries.
type Group struct {
	Name    string
	Entries []diff.Entry
}

// HasDrift reports whether any entry in the group is a drift entry.
func (g Group) HasDrift() bool {
	for _, e := range g.Entries {
		if e.IsDrift() {
			return true
		}
	}
	return false
}

// Count returns the number of entries in the group.
func (g Group) Count() int { return len(g.Entries) }

// Apply groups entries using the provided GroupFunc and returns groups sorted by name.
func Apply(entries []diff.Entry, fn GroupFunc) []Group {
	index := make(map[string][]diff.Entry)
	for _, e := range entries {
		key := fn(e)
		index[key] = append(index[key], e)
	}

	groups := make([]Group, 0, len(index))
	for name, es := range index {
		groups = append(groups, Group{Name: name, Entries: es})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups
}

// ByStatus groups entries by their drift status string.
func ByStatus(entries []diff.Entry) []Group {
	return Apply(entries, func(e diff.Entry) string {
		return string(e.Status)
	})
}

// ByPrefix groups entries by the first segment of the key split on "_".
func ByPrefix(entries []diff.Entry) []Group {
	return Apply(entries, func(e diff.Entry) string {
		for i, c := range e.Key {
			if c == '_' {
				return e.Key[:i]
			}
		}
		return e.Key
	})
}
