// Package ignore provides a mechanism to suppress known or expected drift
// entries from comparison results. Keys can be matched exactly or via
// prefix wildcard patterns (e.g. "DEBUG_*").
//
// This is useful when certain environment variables are intentionally
// different across environments and should not be flagged as drift.
package ignore
