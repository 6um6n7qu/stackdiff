// Package quota provides per-key rate limiting for drift entries.
//
// It tracks how many times each configuration key has drifted within a
// rolling time window and suppresses entries that exceed the configured
// maximum, preventing noisy keys from overwhelming downstream pipelines.
//
// Example usage:
//
//	enforcer := quota.New(quota.DefaultConfig())
//	allowed := enforcer.Filter(entries)
package quota
