// Package flatten provides utilities to flatten nested config maps
// into dot-notation key-value pairs compatible with diff.Entry.
package flatten

import (
	"fmt"
	"sort"

	"github.com/stackdiff/stackdiff/internal/diff"
)

// Flatten converts a nested map[string]any into a flat slice of diff.Entry
// using dot-separated keys. Only leaf string values are included.
func Flatten(prefix string, m map[string]any) []diff.Entry {
	var entries []diff.Entry
	for _, k := range sortedKeys(m) {
		v := m[k]
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]any:
			entries = append(entries, Flatten(fullKey, val)...)
		case string:
			entries = append(entries, diff.Entry{
				Key:    fullKey,
				NewVal: val,
				Status: diff.StatusEqual,
			})
		default:
			entries = append(entries, diff.Entry{
				Key:    fullKey,
				NewVal: fmt.Sprintf("%v", val),
				Status: diff.StatusEqual,
			})
		}
	}
	return entries
}

// ToMap converts a flat slice of diff.Entry back into a map[string]string
// keyed by entry Key using NewVal.
func ToMap(entries []diff.Entry) map[string]string {
	out := make(map[string]string, len(entries))
	for _, e := range entries {
		out[e.Key] = e.NewVal
	}
	return out
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
