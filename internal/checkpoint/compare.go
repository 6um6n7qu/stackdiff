package checkpoint

import (
	"fmt"

	"github.com/user/stackdiff/internal/diff"
)

// Result holds the outcome of comparing a checkpoint against live entries.
type Result struct {
	CheckpointName string
	Drifted        []diff.Entry
	Added          []diff.Entry
	Removed        []diff.Entry
	Changed        []diff.Entry
}

// HasDrift reports whether any drift was detected.
func (r *Result) HasDrift() bool {
	return len(r.Drifted) > 0
}

// Summary returns a human-readable one-line summary.
func (r *Result) Summary() string {
	return fmt.Sprintf("checkpoint %q: +%d -%d ~%d",
		r.CheckpointName, len(r.Added), len(r.Removed), len(r.Changed))
}

// Compare compares a saved checkpoint's entries against live entries and
// returns a Result describing what has changed.
func Compare(cp *Checkpoint, live []diff.Entry) *Result {
	base := indexEntries(cp.Entries)
	current := indexEntries(live)

	res := &Result{CheckpointName: cp.Name}

	for key, liveEntry := range current {
		baseEntry, exists := base[key]
		if !exists {
			e := liveEntry
			e.Status = diff.StatusAdded
			res.Added = append(res.Added, e)
			res.Drifted = append(res.Drifted, e)
			continue
		}
		if baseEntry.NewValue != liveEntry.NewValue {
			e := liveEntry
			e.OldValue = baseEntry.NewValue
			e.Status = diff.StatusChanged
			res.Changed = append(res.Changed, e)
			res.Drifted = append(res.Drifted, e)
		}
	}

	for key, baseEntry := range base {
		if _, exists := current[key]; !exists {
			e := baseEntry
			e.Status = diff.StatusRemoved
			res.Removed = append(res.Removed, e)
			res.Drifted = append(res.Drifted, e)
		}
	}

	return res
}

func indexEntries(entries []diff.Entry) map[string]diff.Entry {
	m := make(map[string]diff.Entry, len(entries))
	for _, e := range entries {
		m[e.Key] = e
	}
	return m
}
