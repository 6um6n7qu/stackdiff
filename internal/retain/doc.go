// Package retain implements a retention policy for diff entries.
//
// Entries can be filtered by status (changed, added, removed) and by
// age using a configurable MaxAge window. Entries without a timestamp
// in their Meta map are always kept when MaxAge is set.
//
// Usage:
//
//	cfg := retain.DefaultConfig()
//	cfg.MaxAge = 24 * time.Hour
//	kept := retain.Apply(entries, cfg)
package retain
