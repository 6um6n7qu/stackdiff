// Package cache provides a simple in-memory TTL cache for diff entries,
// reducing redundant comparisons when configs are polled frequently.
package cache

import (
	"sync"
	"time"

	"github.com/stackdiff/stackdiff/internal/diff"
)

// Entry holds a cached set of diff entries along with expiry metadata.
type Entry struct {
	Entries   []diff.Entry
	CachedAt  time.Time
	ExpiresAt time.Time
}

// IsExpired reports whether the cache entry has passed its TTL.
func (e Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache is a thread-safe in-memory store keyed by an arbitrary string.
type Cache struct {
	mu    sync.RWMutex
	store map[string]Entry
	ttl   time.Duration
}

// DefaultTTL is used when no TTL is specified.
const DefaultTTL = 30 * time.Second

// New creates a Cache with the given TTL. If ttl is zero, DefaultTTL is used.
func New(ttl time.Duration) *Cache {
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	return &Cache{
		store: make(map[string]Entry),
		ttl:   ttl,
	}
}

// Set stores entries under key, overwriting any existing value.
func (c *Cache) Set(key string, entries []diff.Entry) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = Entry{
		Entries:   entries,
		CachedAt:  now,
		ExpiresAt: now.Add(c.ttl),
	}
}

// Get returns the entries for key and whether a valid (non-expired) entry exists.
func (c *Cache) Get(key string) ([]diff.Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.store[key]
	if !ok || e.IsExpired() {
		return nil, false
	}
	return e.Entries, true
}

// Invalidate removes the entry for key if it exists.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]Entry)
}

// Size returns the number of entries currently held (including expired ones).
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}
