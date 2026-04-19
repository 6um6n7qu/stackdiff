package throttle

import (
	"sync"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// Config controls throttle behaviour.
type Config struct {
	// MaxEvents is the maximum number of drift events allowed within Window.
	MaxEvents int
	// Window is the rolling time window.
	Window time.Duration
}

// DefaultConfig returns sensible throttle defaults.
func DefaultConfig() Config {
	return Config{
		MaxEvents: 5,
		Window:    time.Minute,
	}
}

// Throttle tracks event counts over a rolling window.
type Throttle struct {
	cfg    Config
	mu     sync.Mutex
	events []time.Time
}

// New creates a Throttle with the given config.
func New(cfg Config) *Throttle {
	if cfg.MaxEvents <= 0 {
		cfg.MaxEvents = DefaultConfig().MaxEvents
	}
	if cfg.Window <= 0 {
		cfg.Window = DefaultConfig().Window
	}
	return &Throttle{cfg: cfg}
}

// Allow returns true if the event should be allowed through, false if throttled.
func (t *Throttle) Allow(now time.Time) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := now.Add(-t.cfg.Window)
	filtered := t.events[:0]
	for _, ts := range t.events {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}
	t.events = filtered

	if len(t.events) >= t.cfg.MaxEvents {
		return false
	}
	t.events = append(t.events, now)
	return true
}

// Filter returns only the entries that pass the throttle gate.
func (t *Throttle) Filter(entries []diff.Entry, now time.Time) []diff.Entry {
	if !t.Allow(now) {
		return nil
	}
	return entries
}

// Count returns the number of events recorded in the current window.
func (t *Throttle) Count(now time.Time) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	cutoff := now.Add(-t.cfg.Window)
	count := 0
	for _, ts := range t.events {
		if ts.After(cutoff) {
			count++
		}
	}
	return count
}
