// Package digest provides deterministic fingerprinting of diff entry sets.
//
// Use Compute to obtain a SHA-256 hash over a []diff.Entry slice. The hash is
// order-independent — entries are sorted by key before hashing — so two runs
// that produce the same logical drift will always yield the same digest.
//
// Typical usage:
//
//	before := digest.Compute(previousEntries)
//	after  := digest.Compute(currentEntries)
//	if digest.Equal(before, after) {
//	    fmt.Println("no change since last run")
//	}
package digest
