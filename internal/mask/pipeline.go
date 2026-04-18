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

func maskIfNonEmpty(val, replacement string) string {
	if val == "" {
		return val
	}
	return replacement
}
