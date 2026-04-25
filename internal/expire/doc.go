// Package expire implements TTL-based expiration for diff entries.
//
// Entries may carry an "expires_at" metadata field (RFC3339 timestamp).
// When Apply is called, any entry whose expiry time is before the current
// wall-clock time is silently dropped from the result slice.
//
// The expiry key and time format are configurable via Config, and the
// clock source is injectable so that tests can control time without
// sleeping or relying on real wall-clock values.
package expire
