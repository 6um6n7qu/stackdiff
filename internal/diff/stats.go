package diff

// Stats holds aggregate counts derived from a slice of Entry values.
type Stats struct {
	Added   int
	Removed int
	Changed int
	Equal   int
	Total   int
}

// HasDrift reports whether any drift was detected.
func (s Stats) HasDrift() bool {
	return s.Added > 0 || s.Removed > 0 || s.Changed > 0
}

// DriftCount returns the total number of drifted entries.
func (s Stats) DriftCount() int {
	return s.Added + s.Removed + s.Changed
}

// ComputeStats calculates drift statistics from a slice of entries.
func ComputeStats(entries []Entry) Stats {
	var s Stats
	s.Total = len(entries)
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			s.Added++
		case StatusRemoved:
			s.Removed++
		case StatusChanged:
			s.Changed++
		default:
			s.Equal++
		}
	}
	return s
}
