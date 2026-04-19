// Package debounce suppresses repeated drift events within a cooldown window.
package debounce

import (
	"sync"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// Config holds debounce settings.
type Config struct {
	Cooldown time.Duration
}

// DefaultConfig returns a Config with a 30-second cooldown.
func DefaultConfig() Config {
	return Config{Cooldown: 30 * time.Second}
}

// Debouncer tracks last-seen drift fingerprints and suppresses repeats.
type Debouncer struct {
	mu       sync.Mutex
	cooldown time.Duration
	seen     map[string]time.Time
}

// New creates a Debouncer from the given Config.
func New(cfg Config) *Debouncer {
	cd := cfg.Cooldown
	if cd <= 0 {
		cd = DefaultConfig().Cooldown
	}
	return &Debouncer{
		cooldown: cd,
		seen:     make(map[string]time.Time),
	}
}

// Allow returns true if the entry has not been seen within the cooldown window.
// It updates the last-seen timestamp when returning true.
func (d *Debouncer) Allow(e diff.Entry) bool {
	key := e.Key + "|" + string(e.Status)
	d.mu.Lock()
	defer d.mu.Unlock()
	if last, ok := d.seen[key]; ok && time.Since(last) < d.cooldown {
		return false
	}
	d.seen[key] = time.Now()
	return true
}

// Filter returns only entries that pass the debounce check.
func (d *Debouncer) Filter(entries []diff.Entry) []diff.Entry {
	out := make([]diff.Entry, 0, len(entries))
	for _, e := range entries {
		if d.Allow(e) {
			out = append(out, e)
		}
	}
	return out
}

// Reset clears all tracked entries.
func (d *Debouncer) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}
