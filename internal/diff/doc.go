// Package diff provides primitives for comparing two sets of configuration
// entries and representing the resulting differences.
//
// Core types and functions:
//
//   - Entry / Status  – a single key/value pair annotated with its drift status
//     (equal, added, removed, or changed).
//   - Compare          – produces a slice of Entry values from two config maps.
//   - Print            – formats a diff for human-readable output.
//   - Chain            – applies a sequence of transformation steps to entries.
//   - MergeEntries     – combines two entry slices with configurable conflict strategy.
//   - ComputeStats     – aggregates counts per status.
//   - ApplyWindow      – filters entries to a time-based window.
//   - Rank / RankEntries – orders entries by drift severity so the most
//     impactful changes surface first.
package diff
