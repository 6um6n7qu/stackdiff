package history

import (
	"fmt"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// NewRecord builds a Record from a diff result, assigning a timestamp-based ID.
func NewRecord(leftLabel, rightLabel string, entries []diff.Entry) Record {
	now := time.Now().UTC()
	return Record{
		ID:         fmt.Sprintf("%d", now.UnixNano()),
		Timestamp:  now,
		LeftLabel:  leftLabel,
		RightLabel: rightLabel,
		Entries:    entries,
	}
}

// DriftCount returns the number of entries that represent drift (non-equal).
func (r Record) DriftCount() int {
	count := 0
	for _, e := range r.Entries {
		if e.IsDrift() {
			count++
		}
	}
	return count
}

// HasDrift reports whether the record contains any drifted entries.
func (r Record) HasDrift() bool {
	return r.DriftCount() > 0
}

// Summary returns a human-readable one-line summary of the record.
func (r Record) Summary() string {
	return fmt.Sprintf("[%s] %s vs %s — %d drift(s) across %d key(s)",
		r.Timestamp.Format("2006-01-02T15:04:05Z"),
		r.LeftLabel,
		r.RightLabel,
		r.DriftCount(),
		len(r.Entries),
	)
}
