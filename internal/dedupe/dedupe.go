// Package dedupe provides deduplication of diff entries by key,
// keeping the most recently seen entry when duplicates are present.
package dedupe

import "github.com/yourorg/stackdiff/internal/diff"

// Strategy controls which entry wins when duplicates are found.
type Strategy int

const (
	// StrategyFirst keeps the first occurrence of a key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the last occurrence of a key.
	StrategyLast
)

// Apply removes duplicate entries from entries, using the given strategy
// to decide which entry to retain when the same key appears more than once.
// The relative order of retained entries is preserved.
func Apply(entries []diff.Entry, strategy Strategy) []diff.Entry {
	if len(entries) == 0 {
		return entries
	}

	seen := make(map[string]int, len(entries)) // key -> index in result
	result := make([]diff.Entry, 0, len(entries))

	for _, e := range entries {
		if idx, exists := seen[e.Key]; exists {
			if strategy == StrategyLast {
				result[idx] = e
			}
			// StrategyFirst: do nothing
			continue
		}
		seen[e.Key] = len(result)
		result = append(result, e)
	}

	return result
}

// Count returns the number of duplicate keys found in entries.
func Count(entries []diff.Entry) int {
	seen := make(map[string]struct{}, len(entries))
	dupes := 0
	for _, e := range entries {
		if _, exists := seen[e.Key]; exists {
			dupes++
		} else {
			seen[e.Key] = struct{}{}
		}
	}
	return dupes
}
