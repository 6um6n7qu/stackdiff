package env

import "fmt"

// MergeStrategy defines how conflicts are resolved when merging env maps.
type MergeStrategy int

const (
	// StrategyLeft keeps the left-hand value on conflict.
	StrategyLeft MergeStrategy = iota
	// StrategyRight keeps the right-hand value on conflict.
	StrategyRight
	// StrategyError returns an error on conflict.
	StrategyError
)

// MergeResult holds the merged map and metadata about conflicts.
type MergeResult struct {
	Values    map[string]string
	Conflicts []string
}

// MergeMaps merges two env maps using the given strategy.
// Keys present in both maps are handled according to strategy.
func MergeMaps(left, right map[string]string, strategy MergeStrategy) (*MergeResult, error) {
	result := &MergeResult{
		Values:    make(map[string]string),
		Conflicts: []string{},
	}

	for k, v := range left {
		result.Values[k] = v
	}

	for k, v := range right {
		existing, exists := result.Values[k]
		if !exists {
			result.Values[k] = v
			continue
		}
		if existing == v {
			continue
		}
		result.Conflicts = append(result.Conflicts, k)
		switch strategy {
		case StrategyLeft:
			// keep existing
		case StrategyRight:
			result.Values[k] = v
		case StrategyError:
			return nil, fmt.Errorf("conflict on key %q: %q vs %q", k, existing, v)
		}
	}

	return result, nil
}
