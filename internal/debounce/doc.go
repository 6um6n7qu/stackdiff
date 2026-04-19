// Package debounce provides a cooldown-based filter for drift entries.
//
// It prevents the same drift event from being emitted repeatedly within
// a configurable time window, reducing noise in alerting pipelines.
//
// Usage:
//
//	d := debounce.New(debounce.DefaultConfig())
//	filtered := d.Filter(entries)
package debounce
