package snapshot

import (
	"fmt"

	"github.com/stackdiff/internal/diff"
)

// DeltaEntry describes how a single diff.Entry changed between two snapshots.
type DeltaEntry struct {
	Key    string      `json:"key"`
	Before diff.Entry  `json:"before"`
	After  diff.Entry  `json:"after"`
	Change string      `json:"change"` // "added", "removed", "modified", "unchanged"
}

// Delta holds the full comparison result between two snapshots.
type Delta struct {
	Old    *Snapshot    `json:"old"`
	New    *Snapshot    `json:"new"`
	Deltas []DeltaEntry `json:"deltas"`
}

// CompareSnapshots returns a Delta describing what changed between two snapshots.
func CompareSnapshots(old, newer *Snapshot) (*Delta, error) {
	if old == nil || newer == nil {
		return nil, fmt.Errorf("snapshot: cannot compare nil snapshots")
	}

	oldMap := indexEntries(old.Entries)
	newMap := indexEntries(newer.Entries)

	var deltas []DeltaEntry

	for key, newEntry := range newMap {
		if oldEntry, ok := oldMap[key]; !ok {
			deltas = append(deltas, DeltaEntry{Key: key, Before: diff.Entry{}, After: newEntry, Change: "added"})
		} else if oldEntry.ValueA != newEntry.ValueA || oldEntry.ValueB != newEntry.ValueB || oldEntry.Status != newEntry.Status {
			deltas = append(deltas, DeltaEntry{Key: key, Before: oldEntry, After: newEntry, Change: "modified"})
		}
	}

	for key, oldEntry := range oldMap {
		if _, ok := newMap[key]; !ok {
			deltas = append(deltas, DeltaEntry{Key: key, Before: oldEntry, After: diff.Entry{}, Change: "removed"})
		}
	}

	return &Delta{Old: old, New: newer, Deltas: deltas}, nil
}

func indexEntries(entries []diff.Entry) map[string]diff.Entry {
	m := make(map[string]diff.Entry, len(entries))
	for _, e := range entries {
		m[e.Key] = e
	}
	return m
}
