// Package stale provides utilities for identifying drift entries that have
// persisted across multiple snapshots without resolution.
//
// An entry is considered stale when it has appeared as drifted in at least
// MinRepeats consecutive snapshots and its first occurrence exceeds MaxAge.
//
// Typical usage:
//
//	cfg := stale.DefaultConfig()
//	results := stale.Detect(snapshots, timestamps, cfg)
//	for _, r := range results {
//		fmt.Printf("stale key %s (repeats: %d)\n", r.Entry.Key, r.Repeats)
//	}
package stale
