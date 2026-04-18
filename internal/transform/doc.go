// Package transform provides composable transformation functions for diff
// entries. Transformations can normalise keys, trim whitespace, add prefixes,
// or apply any user-supplied mapping before entries are rendered or exported.
//
// Usage:
//
//	t := transform.New(transform.TrimSpace(), transform.LowercaseKeys())
//	result := t.Apply(entries)
package transform
