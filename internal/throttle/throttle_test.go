package throttle_test

import (
	"testing"
	"time"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/throttle"
)

func makeEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "DB_HOST", NewValue: "prod", Status: diff.StatusChanged},
		{Key: "API_KEY", NewValue: "", Status: diff.StatusRemoved},
	}
}

func TestAllow_UnderLimit(t *testing.T) {
	th := throttle.New(throttle.Config{MaxEvents: 3, Window: time.Minute})
	now := time.Now()
	if !th.Allow(now) {
		t.Fatal("expected first event to be allowed")
	}
	if !th.Allow(now.Add(time.Second)) {
		t.Fatal("expected second event to be allowed")
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	th := throttle.New(throttle.Config{MaxEvents: 2, Window: time.Minute})
	now := time.Now()
	th.Allow(now)
	th.Allow(now.Add(time.Second))
	if th.Allow(now.Add(2 * time.Second)) {
		t.Fatal("expected third event to be throttled")
	}
}

func TestAllow_WindowExpiry(t *testing.T) {
	th := throttle.New(throttle.Config{MaxEvents: 1, Window: time.Second})
	now := time.Now()
	th.Allow(now)
	// advance beyond window
	if !th.Allow(now.Add(2 * time.Second)) {
		t.Fatal("expected event after window expiry to be allowed")
	}
}

func TestCount_ReflectsWindow(t *testing.T) {
	th := throttle.New(throttle.Config{MaxEvents: 10, Window: time.Minute})
	now := time.Now()
	th.Allow(now)
	th.Allow(now.Add(time.Second))
	if got := th.Count(now.Add(2 * time.Second)); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}

func TestFilter_AllowedReturnsEntries(t *testing.T) {
	th := throttle.New(throttle.DefaultConfig())
	entries := makeEntries()
	result := th.Filter(entries, time.Now())
	if len(result) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(result))
	}
}

func TestFilter_ThrottledReturnsNil(t *testing.T) {
	th := throttle.New(throttle.Config{MaxEvents: 1, Window: time.Minute})
	now := time.Now()
	th.Allow(now) // consume the slot
	result := th.Filter(makeEntries(), now.Add(time.Second))
	if result != nil {
		t.Fatal("expected nil entries when throttled")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := throttle.DefaultConfig()
	if cfg.MaxEvents <= 0 {
		t.Error("expected positive MaxEvents")
	}
	if cfg.Window <= 0 {
		t.Error("expected positive Window")
	}
}
