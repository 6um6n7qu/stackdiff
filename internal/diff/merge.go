package diff

// MergeStrategy controls how conflicting entries are resolved.
type MergeStrategy int

const (
	// MergeStrategyLeft keeps the value from the left (base) entry.
	MergeStrategyLeft MergeStrategy = iota
	// MergeStrategyRight keeps the value from the right (overlay) entry.
	MergeStrategyRight
	// MergeStrategyUnion includes all entries from both sides.
	MergeStrategyUnion
)

// MergeOptions configures the Merge operation.
type MergeOptions struct {
	Strategy MergeStrategy
}

// DefaultMergeOptions returns sensible defaults.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{Strategy: MergeStrategyRight}
}

// MergeEntries combines two slices of Entry using the given strategy.
// Entries with the same key are resolved according to the strategy.
// Entries unique to either side are always included.
func MergeEntries(base, overlay []Entry, opts MergeOptions) []Entry {
	index := make(map[string]Entry, len(base))
	for _, e := range base {
		index[e.Key] = e
	}

	seen := make(map[string]bool)
	var result []Entry

	for _, e := range overlay {
		seen[e.Key] = true
		if b, ok := index[e.Key]; ok {
			switch opts.Strategy {
			case MergeStrategyLeft:
				result = append(result, b)
			case MergeStrategyRight:
				result = append(result, e)
			case MergeStrategyUnion:
				result = append(result, b)
				if b.Key != e.Key || b.OldVal != e.OldVal || b.NewVal != e.NewVal {
					result = append(result, e)
				}
			}
		} else {
			result = append(result, e)
		}
	}

	for _, e := range base {
		if !seen[e.Key] {
			result = append(result, e)
		}
	}

	return result
}
