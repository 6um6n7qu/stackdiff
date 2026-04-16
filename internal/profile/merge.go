package profile

import "fmt"

// MergeStrategy controls how conflicting profile keys are resolved.
type MergeStrategy int

const (
	StrategyLeft  MergeStrategy = iota // keep base value
	StrategyRight                       // override with overlay value
	StrategyError                       // return error on conflict
)

// MergeResult holds the merged config map and any conflicts encountered.
type MergeResult struct {
	Merged    map[string]string
	Conflicts []string
}

// MergeProfiles merges two profiles by combining their Env maps.
// The base profile is merged with the overlay profile using the given strategy.
func MergeProfiles(base, overlay Profile, strategy MergeStrategy) (MergeResult, error) {
	result := MergeResult{
		Merged: make(map[string]string),
	}

	for k, v := range base.Env {
		result.Merged[k] = v
	}

	for k, v := range overlay.Env {
		existing, exists := result.Merged[k]
		if !exists {
			result.Merged[k] = v
			continue
		}
		if existing == v {
			continue
		}
		switch strategy {
		case StrategyLeft:
			// keep existing, do nothing
		case StrategyRight:
			result.Merged[k] = v
		case StrategyError:
			return MergeResult{}, fmt.Errorf("conflict on key %q: %q vs %q", k, existing, v)
		}
		result.Conflicts = append(result.Conflicts, k)
	}

	return result, nil
}
