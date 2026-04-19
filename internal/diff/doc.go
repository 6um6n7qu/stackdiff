// Package diff provides types and utilities for comparing two sets of
// configuration entries and representing the resulting drift.
//
// Core concepts:
//
//   - Entry: a single key/value pair with an associated drift status.
//   - Compare: compares two maps and returns a slice of Entry values.
//   - Chain: applies an ordered sequence of named transformations to entries,
//     optionally recording intermediate results for tracing.
//   - Print: renders a slice of entries to an io.Writer.
package diff
