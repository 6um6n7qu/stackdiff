// Package replay walks a history store and re-emits records through a
// caller-supplied handler, optionally filtered by time window or drift
// presence. It is useful for re-processing past comparisons without
// re-running the original diff.
package replay
