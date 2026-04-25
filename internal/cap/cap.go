// Package cap limits the number of drift entries returned per status bucket.
// It is useful when downstream consumers need a bounded result set.
package cap

import "github.com/user/stackdiff/internal/diff"

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxPerStatus: 10,
		SentinelKey:  "__cap_truncated__",
	}
}

// Config controls how Apply truncates entries.
type Config struct {
	// MaxPerStatus is the maximum number of entries allowed per status bucket.
	// Zero means no limit.
	MaxPerStatus int

	// SentinelKey is the key used for a synthetic entry appended when a bucket
	// is truncated, so callers can detect that capping occurred.
	SentinelKey string
}

// Apply returns a new slice capped to Config.MaxPerStatus entries per status.
// When a bucket is truncated a sentinel entry is appended with StatusChanged
// and a value describing how many entries were dropped.
func Apply(entries []diff.Entry, cfg Config) []diff.Entry {
	if cfg.MaxPerStatus <= 0 {
		return entries
	}
	if cfg.SentinelKey == "" {
		cfg.SentinelKey = DefaultConfig().SentinelKey
	}

	counts := make(map[diff.Status]int)
	result := make([]diff.Entry, 0, len(entries))

	for _, e := range entries {
		if counts[e.Status] < cfg.MaxPerStatus {
			result = append(result, e)
			counts[e.Status]++
		}
	}

	// Append sentinels for each truncated bucket.
	for _, e := range entries {
		total := countStatus(entries, e.Status)
		if total > cfg.MaxPerStatus {
			// Only add one sentinel per status.
			if !hasSentinel(result, cfg.SentinelKey, e.Status) {
				result = append(result, diff.Entry{
					Key:      cfg.SentinelKey,
					NewValue: fmt.Sprintf("%d entries truncated", total-cfg.MaxPerStatus),
					Status:   e.Status,
				})
			}
		}
	}
	return result
}

// Truncated returns true when Apply would truncate at least one bucket.
func Truncated(entries []diff.Entry, cfg Config) bool {
	if cfg.MaxPerStatus <= 0 {
		return false
	}
	counts := make(map[diff.Status]int)
	for _, e := range entries {
		counts[e.Status]++
		if counts[e.Status] > cfg.MaxPerStatus {
			return true
		}
	}
	return false
}

func countStatus(entries []diff.Entry, s diff.Status) int {
	n := 0
	for _, e := range entries {
		if e.Status == s {
			n++
		}
	}
	return n
}

func hasSentinel(entries []diff.Entry, key string, s diff.Status) bool {
	for _, e := range entries {
		if e.Key == key && e.Status == s {
			return true
		}
	}
	return false
}
