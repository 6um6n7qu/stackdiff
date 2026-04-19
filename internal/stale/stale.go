// Package stale identifies drift entries that have not changed across
// multiple consecutive snapshots and flags them as stale.
package stale

import (
	"time"

	"stackdiff/internal/diff"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxAge:     72 * time.Hour,
		MinRepeats: 3,
	}
}

// Config controls how staleness is determined.
type Config struct {
	// MaxAge is the maximum duration an entry can persist without change.
	MaxAge time.Duration
	// MinRepeats is the minimum number of times an entry must appear unchanged.
	MinRepeats int
}

// Result holds a stale entry alongside its repeat count.
type Result struct {
	Entry   diff.Entry
	Repeats int
	Since   time.Time
}

// Detect returns entries considered stale given a series of entry slices
// ordered oldest-first. An entry is stale if it appears unchanged in at
// least cfg.MinRepeats consecutive snapshots and its first occurrence is
// older than cfg.MaxAge relative to now.
func Detect(snapshots [][]diff.Entry, timestamps []time.Time, cfg Config) []Result {
	if len(snapshots) == 0 {
		return nil
	}
	if cfg.MinRepeats <= 0 {
		cfg.MinRepeats = DefaultConfig().MinRepeats
	}
	if cfg.MaxAge <= 0 {
		cfg.MaxAge = DefaultConfig().MaxAge
	}

	// Track first-seen time and repeat count per key.
	type meta struct {
		firstSeen time.Time
		repeats   int
		entry     diff.Entry
	}
	tracked := map[string]*meta{}

	for i, snap := range snapshots {
		var ts time.Time
		if i < len(timestamps) {
			ts = timestamps[i]
		} else {
			ts = time.Now()
		}
		for _, e := range snap {
			if !e.IsDrift() {
				continue
			}
			if m, ok := tracked[e.Key]; ok && m.entry.NewValue == e.NewValue {
				m.repeats++
				m.entry = e
			} else {
				tracked[e.Key] = &meta{firstSeen: ts, repeats: 1, entry: e}
			}
		}
	}

	now := time.Now()
	var results []Result
	for _, m := range tracked {
		age := now.Sub(m.firstSeen)
		if m.repeats >= cfg.MinRepeats && age >= cfg.MaxAge {
			results = append(results, Result{
				Entry:   m.entry,
				Repeats: m.repeats,
				Since:   m.firstSeen,
			})
		}
	}
	return results
}
