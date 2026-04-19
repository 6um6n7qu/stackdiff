package mask

import "github.com/user/stackdiff/internal/diff"

// ApplyToEntries returns a copy of entries with sensitive values masked.
// Both OldVal and NewVal are masked for matching keys.
func ApplyToEntries(m *Masker, entries []diff.Entry) []diff.Entry {
	result := make([]diff.Entry, len(entries))
	for i, e := range entries {
		if m.IsSensitive(e.Key) {
			replacement := m.replacementFor(e.Key)
			e.OldVal = maskIfNonEmpty(e.OldVal, replacement)
			e.NewVal = maskIfNonEmpty(e.NewVal, replacement)
		}
		result[i] = e
	}
	return result
}

// CountMasked returns the number of entries whose values would be masked.
func CountMasked(m *Masker, entries []diff.Entry) int {
	count := 0
	for _, e := range entries {
		if m.IsSensitive(e.Key) {
			count++
		}
	}
	return count
}

func maskIfNonEmpty(val, replacement string) string {
	if val == "" {
		return val
	}
	return replacement
}
