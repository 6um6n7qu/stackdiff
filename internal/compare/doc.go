// Package compare provides a high-level pipeline that wires together diff,
// normalization, masking, filtering, and enrichment into a single call.
//
// Basic usage:
//
//	r, err := compare.Run(mapA, mapB, compare.DefaultOptions())
//	if err != nil { ... }
//	if r.HasDrift() { ... }
package compare
