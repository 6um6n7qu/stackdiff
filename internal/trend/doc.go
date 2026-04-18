// Package trend provides drift frequency analysis over historical records.
//
// Use Analyze to scan a slice of history.Record values within a time window
// and identify config keys that have drifted repeatedly. Results are sorted
// by frequency descending so the most volatile keys appear first.
package trend
