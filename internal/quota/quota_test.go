package quota_test

import (
	"testing"
	"time"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/quota"
)

func makeEntry(key string) diff.Entry {
	return diff.Entry{Key: key, NewValue: "v", Status: diff.StatusChanged}
}

func TestAllow_UnderLimit(t *testing.T) {
	e := quota.New(quota.Config{MaxPerKey: 3, Window: time.Minute})
	entry := makeEntry("DB_HOST")
	for i := 0; i < 3; i++ {
		if !e.Allow(entry) {
			t.Fatalf("expected Allow=true on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	e := quota.New(quota.Config{MaxPerKey: 2, Window: time.Minute})
	entry := makeEntry("DB_HOST")
	e.Allow(entry)
	e.Allow(entry)
	if e.Allow(entry) {
		t.Fatal("expected Allow=false after exceeding quota")
	}
}

func TestAllow_WindowExpiry(t *testing.T) {
	now := time.Now()
	e := quota.New(quota.Config{MaxPerKey: 1, Window: time.Millisecond})
	// inject controllable clock
	clock := now
	e2 := &struct{ *quota.Enforcer }{quota.New(quota.Config{MaxPerKey: 1, Window: time.Millisecond})}
	_ = e2
	_ = clock

	// Use the exported path: after window expires count resets
	entry := makeEntry("KEY")
	e.Allow(entry) // count=1, at limit
	time.Sleep(5 * time.Millisecond)
	if !e.Allow(entry) {
		t.Fatal("expected Allow=true after window expiry")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	e := quota.New(quota.Config{MaxPerKey: 1, Window: time.Minute})
	a := makeEntry("A")
	b := makeEntry("B")
	if !e.Allow(a) {
		t.Fatal("A should be allowed")
	}
	if !e.Allow(b) {
		t.Fatal("B should be allowed independently")
	}
	if e.Allow(a) {
		t.Fatal("A should be blocked after quota")
	}
}

func TestFilter_RemovesOverQuota(t *testing.T) {
	e := quota.New(quota.Config{MaxPerKey: 2, Window: time.Minute})
	entries := []diff.Entry{
		makeEntry("X"), makeEntry("X"), makeEntry("X"),
		makeEntry("Y"),
	}
	got := e.Filter(entries)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestCount_ReflectsWindow(t *testing.T) {
	e := quota.New(quota.Config{MaxPerKey: 10, Window: time.Minute})
	entry := makeEntry("Z")
	e.Allow(entry)
	e.Allow(entry)
	if c := e.Count("Z"); c != 2 {
		t.Fatalf("expected count 2, got %d", c)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := quota.DefaultConfig()
	if cfg.MaxPerKey <= 0 {
		t.Error("MaxPerKey should be positive")
	}
	if cfg.Window <= 0 {
		t.Error("Window should be positive")
	}
}
