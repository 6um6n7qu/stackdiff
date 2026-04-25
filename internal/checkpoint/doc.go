// Package checkpoint provides named, timestamped snapshots of config
// entries that act as reference points for drift detection.
//
// A checkpoint is saved to disk as a JSON file and can later be loaded
// and compared against the current live config to produce a Result
// describing what has been added, removed, or changed since the
// checkpoint was taken.
//
// Typical usage:
//
//	cp := checkpoint.New("pre-deploy", entries)
//	_ = checkpoint.Save(".stackdiff/checkpoints", "pre-deploy", cp)
//	// later ...
//	cp, _ = checkpoint.Load(".stackdiff/checkpoints", "pre-deploy")
//	result := checkpoint.Compare(cp, liveEntries)
//	if result.HasDrift() { ... }
package checkpoint
