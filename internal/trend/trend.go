// Package trend analyzes drift history to surface recurring and worsening config issues.
package trend

import (
	"sort"
	"time"

	"github.com/stackdiff/internal/history"
)

// Entry represents a single key's drift frequency over a time window.
type Entry struct {
	Key       string
	Count     int
	LastSeen  time.Time
	Statuses  []string
}

// Report holds the trend analysis result.
type Report struct {
	Window  time.Duration
	Entries []Entry
}

// Analyze scans history records within the given time window and returns
// a Report of keys that drifted more than once.
func Analyze(records []history.Record, window time.Duration) Report {
	cutoff := time.Now().Add(-window)
	freq := map[string]*Entry{}

	for _, r := range records {
		if r.Timestamp.Before(cutoff) {
			continue
		}
		for _, e := range r.Diffs {
			if !e.IsDrift() {
				continue
			}
			ent, ok := freq[e.Key]
			if !ok {
				ent = &Entry{Key: e.Key}
				freq[e.Key] = ent
			}
			ent.Count++
			if r.Timestamp.After(ent.LastSeen) {
				ent.LastSeen = r.Timestamp
			}
			ent.Statuses = append(ent.Statuses, string(e.Status))
		}
	}

	var entries []Entry
	for _, e := range freq {
		if e.Count > 1 {
			entries = append(entries, *e)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Count > entries[j].Count
	})

	return Report{Window: window, Entries: entries}
}
