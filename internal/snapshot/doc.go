// Package snapshot provides functionality for capturing, persisting, and
// comparing diff results over time. A Snapshot records the output of a
// stackdiff comparison at a specific point in time, including labels for
// the two compared services and all resulting diff entries.
//
// Typical usage:
//
//	// Save a snapshot after running a diff
//	s := snapshot.New("staging", "prod", entries)
//	if err := snapshot.Save(s, "./snapshots/2024-01-15.json"); err != nil {
//		log.Fatal(err)
//	}
//
//	// Load a previous snapshot and compare
//	old, err := snapshot.Load("./snapshots/2024-01-14.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	delta, err := snapshot.CompareSnapshots(old, s)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, d := range delta.Deltas {
//		fmt.Printf("%s: %s\n", d.Key, d.Change)
//	}
package snapshot
