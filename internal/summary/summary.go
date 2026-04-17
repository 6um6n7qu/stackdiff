package summary

import (
	"fmt"
	"strings"

	"github.com/yourorg/stackdiff/internal/diff"
)

// Stats holds aggregated counts of a diff result.
type Stats struct {
	Added   int
	Removed int
	Changed int
	Total   int
}

// Build computes Stats from a slice of diff entries.
func Build(entries []diff.Entry) Stats {
	s := Stats{Total: len(entries)}
	for _, e := range entries {
		switch e.Status {
		case diff.StatusAdded:
			s.Added++
		case diff.StatusRemoved:
			s.Removed++
		case diff.StatusChanged:
			s.Changed++
		}
	}
	return s
}

// HasDrift returns true when any drift was detected.
func (s Stats) HasDrift() bool {
	return s.Added > 0 || s.Removed > 0 || s.Changed > 0
}

// String returns a human-readable one-line summary.
func (s Stats) String() string {
	if !s.HasDrift() {
		return "no drift detected"
	}
	parts := []string{}
	if s.Added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", s.Added))
	}
	if s.Removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", s.Removed))
	}
	if s.Changed > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", s.Changed))
	}
	return fmt.Sprintf("drift detected: %s (total keys: %d)", strings.Join(parts, ", "), s.Total)
}
