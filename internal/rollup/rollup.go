// Package rollup aggregates multiple diff entries into a grouped summary
// keyed by a user-defined dimension (e.g. service, namespace, prefix).
package rollup

import (
	"sort"
	"strings"

	"github.com/yourorg/stackdiff/internal/diff"
)

// Group holds aggregated drift entries for a single dimension value.
type Group struct {
	Key     string
	Entries []diff.Entry
	Added   int
	Removed int
	Changed int
}

// HasDrift returns true when the group contains any drift.
func (g Group) HasDrift() bool {
	return g.Added+g.Removed+g.Changed > 0
}

// KeyFunc extracts a grouping key from an entry key string.
type KeyFunc func(entryKey string) string

// PrefixKeyFunc groups by the first segment of a dot-separated key.
func PrefixKeyFunc(entryKey string) string {
	parts := strings.SplitN(entryKey, ".", 2)
	return parts[0]
}

// ByKeyFunc rolls up entries into groups using the provided KeyFunc.
func ByKeyFunc(entries []diff.Entry, fn KeyFunc) []Group {
	index := map[string]*Group{}

	for _, e := range entries {
		k := fn(e.Key)
		if k == "" {
			k = "(unknown)"
		}
		g, ok := index[k]
		if !ok {
			g = &Group{Key: k}
			index[k] = g
		}
		g.Entries = append(g.Entries, e)
		switch e.Status {
		case diff.StatusAdded:
			g.Added++
		case diff.StatusRemoved:
			g.Removed++
		case diff.StatusChanged:
			g.Changed++
		}
	}

	groups := make([]Group, 0, len(index))
	for _, g := range index {
		groups = append(groups, *g)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Key < groups[j].Key
	})
	return groups
}
