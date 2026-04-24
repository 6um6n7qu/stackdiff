// Package quota enforces per-key drift event limits over a rolling window.
// It prevents a single noisy key from flooding downstream consumers.
package quota

import (
	"sync"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// DefaultConfig returns a sensible quota configuration.
func DefaultConfig() Config {
	return Config{
		MaxPerKey: 5,
		Window:    time.Minute,
	}
}

// Config controls quota enforcement behaviour.
type Config struct {
	MaxPerKey int           // maximum drift events allowed per key in Window
	Window    time.Duration // rolling window duration
}

type bucket struct {
	count int
	reset time.Time
}

// Enforcer tracks per-key event counts and blocks keys that exceed the quota.
type Enforcer struct {
	mu      sync.Mutex
	cfg     Config
	buckets map[string]*bucket
	now     func() time.Time // injectable for testing
}

// New creates a new Enforcer with the given configuration.
func New(cfg Config) *Enforcer {
	if cfg.MaxPerKey <= 0 {
		cfg.MaxPerKey = DefaultConfig().MaxPerKey
	}
	if cfg.Window <= 0 {
		cfg.Window = DefaultConfig().Window
	}
	return &Enforcer{
		cfg:     cfg,
		buckets: make(map[string]*bucket),
		now:     time.Now,
	}
}

// Allow reports whether the given entry is within quota.
// It increments the counter for the entry's key and returns false if the
// quota has been exceeded for the current window.
func (e *Enforcer) Allow(entry diff.Entry) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := e.now()
	b, ok := e.buckets[entry.Key]
	if !ok || now.After(b.reset) {
		e.buckets[entry.Key] = &bucket{count: 1, reset: now.Add(e.cfg.Window)}
		return true
	}
	b.count++
	return b.count <= e.cfg.MaxPerKey
}

// Filter returns only the entries that are within quota.
func (e *Enforcer) Filter(entries []diff.Entry) []diff.Entry {
	out := make([]diff.Entry, 0, len(entries))
	for _, en := range entries {
		if e.Allow(en) {
			out = append(out, en)
		}
	}
	return out
}

// Count returns the current event count for a key (0 if unseen or expired).
func (e *Enforcer) Count(key string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	b, ok := e.buckets[key]
	if !ok || e.now().After(b.reset) {
		return 0
	}
	return b.count
}
