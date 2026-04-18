package clamp_test

import (
	"strings"
	"testing"

	"github.com/user/stackdiff/internal/clamp"
	"github.com/user/stackdiff/internal/diff"
)

func makeEntry(key, oldVal, newVal string, status diff.Status) diff.Entry {
	return diff.Entry{Key: key, OldValue: oldVal, NewValue: newVal, Status: status}
}

func TestClamp_ShortValuesUnchanged(t *testing.T) {
	cfg := clamp.DefaultConfig()
	entries := []diff.Entry{
		makeEntry("KEY", "short", "also_short", diff.StatusChanged),
	}
	out := clamp.Clamp(entries, cfg)
	if out[0].OldValue != "short" || out[0].NewValue != "also_short" {
		t.Errorf("expected values unchanged, got %+v", out[0])
	}
}

func TestClamp_LongValueTruncated(t *testing.T) {
	cfg := clamp.Config{MaxLength: 10, Suffix: "..."}
	long := strings.Repeat("x", 50)
	entries := []diff.Entry{
		makeEntry("BIG", long, long, diff.StatusChanged),
	}
	out := clamp.Clamp(entries, cfg)
	expected := strings.Repeat("x", 10) + "..."
	if out[0].OldValue != expected {
		t.Errorf("expected %q, got %q", expected, out[0].OldValue)
	}
	if out[0].NewValue != expected {
		t.Errorf("expected %q, got %q", expected, out[0].NewValue)
	}
}

func TestClamp_DoesNotMutateOriginal(t *testing.T) {
	cfg := clamp.Config{MaxLength: 5, Suffix: "!"}
	original := strings.Repeat("a", 20)
	entries := []diff.Entry{
		makeEntry("K", original, "", diff.StatusRemoved),
	}
	clamp.Clamp(entries, cfg)
	if entries[0].OldValue != original {
		t.Error("original entry was mutated")
	}
}

func TestClamp_ZeroMaxLengthUsesDefault(t *testing.T) {
	cfg := clamp.Config{MaxLength: 0, Suffix: "..."}
	short := "hello"
	entries := []diff.Entry{
		makeEntry("K", short, short, diff.StatusEqual),
	}
	out := clamp.Clamp(entries, cfg)
	if out[0].OldValue != short {
		t.Errorf("expected %q, got %q", short, out[0].OldValue)
	}
}

func TestClamp_ExactLengthNotTruncated(t *testing.T) {
	cfg := clamp.Config{MaxLength: 5, Suffix: "..."}
	val := "hello"
	entries := []diff.Entry{
		makeEntry("K", val, val, diff.StatusEqual),
	}
	out := clamp.Clamp(entries, cfg)
	if out[0].OldValue != val {
		t.Errorf("expected exact-length value unchanged, got %q", out[0].OldValue)
	}
}
