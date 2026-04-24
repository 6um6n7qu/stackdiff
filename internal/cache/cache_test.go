package cache_test

import (
	"testing"
	"time"

	"github.com/stackdiff/stackdiff/internal/cache"
	"github.com/stackdiff/stackdiff/internal/diff"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "APP_ENV", OldValue: "staging", NewValue: "production", Status: diff.StatusChanged},
		{Key: "LOG_LEVEL", OldValue: "debug", NewValue: "debug", Status: diff.StatusEqual},
	}
}

func TestNew_DefaultTTL(t *testing.T) {
	c := cache.New(0)
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestSetAndGet_HitBeforeExpiry(t *testing.T) {
	c := cache.New(5 * time.Second)
	entries := sampleEntries()
	c.Set("key1", entries)

	got, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(got) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(got))
	}
}

func TestGet_MissOnUnknownKey(t *testing.T) {
	c := cache.New(5 * time.Second)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss for unknown key")
	}
}

func TestGet_MissAfterExpiry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("key1", sampleEntries())
	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("key1")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("key1", sampleEntries())
	c.Invalidate("key1")

	_, ok := c.Get("key1")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("a", sampleEntries())
	c.Set("b", sampleEntries())
	c.Flush()

	if c.Size() != 0 {
		t.Fatalf("expected size 0 after flush, got %d", c.Size())
	}
}

func TestSize_ReflectsEntryCount(t *testing.T) {
	c := cache.New(5 * time.Second)
	if c.Size() != 0 {
		t.Fatalf("expected initial size 0, got %d", c.Size())
	}
	c.Set("x", sampleEntries())
	c.Set("y", sampleEntries())
	if c.Size() != 2 {
		t.Fatalf("expected size 2, got %d", c.Size())
	}
}

func TestSet_OverwritesExistingKey(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("key1", sampleEntries())
	newEntries := []diff.Entry{{Key: "NEW_KEY", NewValue: "val", Status: diff.StatusAdded}}
	c.Set("key1", newEntries)

	got, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected cache hit after overwrite")
	}
	if len(got) != 1 || got[0].Key != "NEW_KEY" {
		t.Fatalf("expected overwritten entry, got %+v", got)
	}
}
