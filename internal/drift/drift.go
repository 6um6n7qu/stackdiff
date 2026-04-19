// Package drift provides a high-level summary of drift severity across entries.
package drift

import (
	"fmt"

	"github.com/yourorg/stackdiff/internal/diff"
)

// Level represents the overall drift severity.
type Level string

const (
	LevelNone     Level = "none"
	LevelLow      Level = "low"
	LevelModerate Level = "moderate"
	LevelHigh     Level = "high"
)

// Result holds the evaluated drift level and contributing counts.
type Result struct {
	Level    Level
	Added    int
	Removed  int
	Changed  int
	Total    int
}

// String returns a human-readable summary.
func (r Result) String() string {
	return fmt.Sprintf("level=%s added=%d removed=%d changed=%d total=%d",
		r.Level, r.Added, r.Removed, r.Changed, r.Total)
}

// Evaluate inspects a slice of diff entries and returns a Result.
func Evaluate(entries []diff.Entry) Result {
	var added, removed, changed int
	for _, e := range entries {
		switch e.Status {
		case diff.StatusAdded:
			added++
		case diff.StatusRemoved:
			removed++
		case diff.StatusChanged:
			changed++
		}
	}
	total := added + removed + changed
	return Result{
		Level:   classify(total, removed),
		Added:   added,
		Removed: removed,
		Changed: changed,
		Total:   total,
	}
}

func classify(total, removed int) Level {
	switch {
	case total == 0:
		return LevelNone
	case removed >= 3 || total >= 10:
		return LevelHigh
	case total >= 4:
		return LevelModerate
	default:
		return LevelLow
	}
}
