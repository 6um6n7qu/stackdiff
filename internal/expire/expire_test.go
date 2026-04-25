package expire

import (
	"testing"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

func makeEntry(key, val, expiresAt string) diff.Entry {
	e := diff.Entry{
		Key:      key,
		NewValue: val,
		Status:   diff.StatusChanged,
	}
	if expiresAt != "" {
		e.Meta = map[string]string{"expires_at": expiresAt}
	}
	return e
}

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func cfgAt(t time.Time) Config {
	c := DefaultConfig()
	c.Now = func() time.Time { return t }
	return c
}

func TestApply_NoExpiry_KeepsAll(t *testing.T) {
	entries := []diff.Entry{makeEntry("A", "1", ""), makeEntry("B", "2", "")}
	out := Apply(entries, cfgAt(fixedNow))
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestApply_FutureExpiry_Kept(t *testing.T) {
	future := fixedNow.Add(time.Hour).Format(time.RFC3339)
	entries := []diff.Entry{makeEntry("X", "v", future)}
	out := Apply(entries, cfgAt(fixedNow))
	if len(out) != 1 {
		t.Fatalf("expected entry to be kept, got %d", len(out))
	}
}

func TestApply_PastExpiry_Removed(t *testing.T) {
	past := fixedNow.Add(-time.Hour).Format(time.RFC3339)
	entries := []diff.Entry{makeEntry("X", "v", past)}
	out := Apply(entries, cfgAt(fixedNow))
	if len(out) != 0 {
		t.Fatalf("expected entry to be removed, got %d", len(out))
	}
}

func TestApply_MixedExpiry(t *testing.T) {
	past := fixedNow.Add(-time.Minute).Format(time.RFC3339)
	future := fixedNow.Add(time.Minute).Format(time.RFC3339)
	entries := []diff.Entry{
		makeEntry("keep", "1", future),
		makeEntry("drop", "2", past),
		makeEntry("none", "3", ""),
	}
	out := Apply(entries, cfgAt(fixedNow))
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestCountExpired(t *testing.T) {
	past := fixedNow.Add(-time.Hour).Format(time.RFC3339)
	entries := []diff.Entry{
		makeEntry("a", "1", past),
		makeEntry("b", "2", past),
		makeEntry("c", "3", ""),
	}
	n := CountExpired(entries, cfgAt(fixedNow))
	if n != 2 {
		t.Fatalf("expected 2 expired, got %d", n)
	}
}

func TestApply_InvalidTimestamp_Kept(t *testing.T) {
	e := diff.Entry{Key: "k", Meta: map[string]string{"expires_at": "not-a-date"}}
	out := Apply([]diff.Entry{e}, cfgAt(fixedNow))
	if len(out) != 1 {
		t.Fatal("entry with unparseable timestamp should be kept")
	}
}
