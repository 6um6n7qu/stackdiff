package debounce_test

import (
	"testing"
	"time"

	"github.com/user/stackdiff/internal/debounce"
	"github.com/user/stackdiff/internal/diff"
)

func makeEntry(key string, status diff.Status) diff.Entry {
	return diff.Entry{Key: key, NewValue: "v", Status: status}
}

func TestAllow_FirstCallReturnsTrue(t *testing.T) {
	d := debounce.New(debounce.DefaultConfig())
	e := makeEntry("PORT", diff.StatusChanged)
	if !d.Allow(e) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldownReturnsFalse(t *testing.T) {
	d := debounce.New(debounce.Config{Cooldown: 10 * time.Second})
	e := makeEntry("PORT", diff.StatusChanged)
	d.Allow(e)
	if d.Allow(e) {
		t.Fatal("expected second call within cooldown to be suppressed")
	}
}

func TestAllow_DifferentStatusTreatedSeparately(t *testing.T) {
	d := debounce.New(debounce.Config{Cooldown: 10 * time.Second})
	d.Allow(makeEntry("PORT", diff.StatusChanged))
	if !d.Allow(makeEntry("PORT", diff.StatusAdded)) {
		t.Fatal("different status for same key should be allowed")
	}
}

func TestFilter_RemovesDuplicates(t *testing.T) {
	d := debounce.New(debounce.Config{Cooldown: 10 * time.Second})
	entries := []diff.Entry{
		makeEntry("A", diff.StatusChanged),
		makeEntry("A", diff.StatusChanged),
		makeEntry("B", diff.StatusAdded),
	}
	out := d.Filter(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := debounce.New(debounce.Config{Cooldown: 10 * time.Second})
	e := makeEntry("X", diff.StatusRemoved)
	d.Allow(e)
	d.Reset()
	if !d.Allow(e) {
		t.Fatal("expected entry to be allowed after reset")
	}
}

func TestDefaultConfig_HasPositiveCooldown(t *testing.T) {
	cfg := debounce.DefaultConfig()
	if cfg.Cooldown <= 0 {
		t.Fatal("default cooldown should be positive")
	}
}
