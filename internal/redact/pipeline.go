package redact

import "github.com/user/stackdiff/internal/diff"

// ApplyToEntries returns a new slice of entries with sensitive values redacted.
// Original entries are not mutated.
func ApplyToEntries(r *Redactor, entries []diff.Entry) []diff.Entry {
	result := make([]diff.Entry, len(entries))
	for i, e := range entries {
		result[i] = diff.Entry{
			Key:    e.Key,
			OldVal: redactIfSensitive(r, e.Key, e.OldVal),
			NewVal: redactIfSensitive(r, e.Key, e.NewVal),
			Status: e.Status,
		}
	}
	return result
}

// ApplyToEntry returns a new entry with sensitive values redacted.
// The original entry is not mutated.
func ApplyToEntry(r *Redactor, e diff.Entry) diff.Entry {
	return diff.Entry{
		Key:    e.Key,
		OldVal: redactIfSensitive(r, e.Key, e.OldVal),
		NewVal: redactIfSensitive(r, e.Key, e.NewVal),
		Status: e.Status,
	}
}

func redactIfSensitive(r *Redactor, key, val string) string {
	if val == "" {
		return val
	}
	if r.IsSensitive(key) {
		return r.Apply(key, val)
	}
	return val
}
