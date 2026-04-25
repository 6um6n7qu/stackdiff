// Package reorder provides sorting and reordering utilities for slices of
// diff.Entry values.
//
// Supported strategies:
//
//   - ByKey      – alphabetical by key name (case-insensitive)
//   - ByKeyDesc  – reverse-alphabetical by key name
//   - ByStatus   – drift-severity order: changed → added → removed → equal
//
// The Apply function never mutates the input slice.
package reorder
