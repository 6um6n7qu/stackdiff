package diff

import "github.com/stackdiff/stackdiff/internal/config"

// DriftKind categorises the type of difference found between two configs.
type DriftKind string

const (
	Added    DriftKind = "added"    // key exists in B but not A
	Removed  DriftKind = "removed"  // key exists in A but not B
	Changed  DriftKind = "changed"  // key exists in both but values differ
)

// DriftEntry represents a single divergence between two configs.
type DriftEntry struct {
	Key      string
	Kind     DriftKind
	ValueA   string
	ValueB   string
}

// Compare returns a slice of DriftEntry describing every difference
// between cfgA and cfgB.
func Compare(cfgA, cfgB config.Config) []DriftEntry {
	var entries []DriftEntry

	// Keys in A — check for removals and changes.
	for k, va := range cfgA {
		if vb, ok := cfgB[k]; !ok {
			entries = append(entries, DriftEntry{Key: k, Kind: Removed, ValueA: va})
		} else if va != vb {
			entries = append(entries, DriftEntry{Key: k, Kind: Changed, ValueA: va, ValueB: vb})
		}
	}

	// Keys only in B — additions.
	for k, vb := range cfgB {
		if _, ok := cfgA[k]; !ok {
			entries = append(entries, DriftEntry{Key: k, Kind: Added, ValueB: vb})
		}
	}

	return entries
}
