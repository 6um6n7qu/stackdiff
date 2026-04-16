// Package alert provides lightweight drift alerting for stackdiff.
//
// When a comparison produces drift entries, Emit can be called to
// surface a human-readable alert to any io.Writer (default: stderr).
//
// Example:
//
//	cfg := alert.DefaultConfig()
//	alert.Emit(driftEntries, cfg)
package alert
