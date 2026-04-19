// Package prune provides utilities for removing stale or aged-out diff entries
// from a working set. It is useful when operating in watch or schedule mode
// where entries accumulate over time and old drift signals should be expired.
//
// Usage:
//
//	cfg := prune.DefaultConfig()
//	cfg.MaxAge = 48 * time.Hour
//	result := prune.Apply(entries, cfg, time.Now())
//	// result.Kept  — entries still within the age window
//	// result.Pruned — entries that were expired
package prune
