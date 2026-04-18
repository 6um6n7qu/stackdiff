package truncate_test

import (
	"strings"
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/truncate"
)

func makeEntry(key, oldVal, newVal string, status diff.Status) diff.Entry {
	return diff.Entry{Key: key, OldVal: oldVal, NewVal: newVal, Status: status}
}

func TestApply_ShortValuesUnchanged(t *testing.T) {
	tr := truncate.New(truncate.DefaultConfig())
	entries := []diff.Entry{
		makeEntry("KEY", "short", "also short", diff.StatusChanged),
	}
	out := tr.Apply(entries)
	if out[0].OldVal != "short" || out[0].NewVal != "also short" {
		t.Errorf("expected values unchanged, got %q %q", out[0].OldVal, out[0].NewVal)
	}
}

func TestApply_LongValueTruncated(t *testing.T) {
	cfg := truncate.Config{MaxLen: 10, Suffix: "..."}
	tr := truncate.New(cfg)
	long := strings.Repeat("x", 50)
	entries := []diff.Entry{
		makeEntry("K", long, long, diff.StatusChanged),
	}
	out := tr.Apply(entries)
	if len(out[0].OldVal) != 13 {
		t.Errorf("expected truncated len 13, got %d", len(out[0].OldVal))
	}
	if !strings.HasSuffix(out[0].OldVal, "...") {
		t.Errorf("expected suffix '...', got %q", out[0].OldVal)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	tr := truncate.New(truncate.Config{MaxLen: 5, Suffix: "~"})
	origVal := strings.Repeat("a", 20)
	entries := []diff.Entry{
		makeEntry("K", origVal, origVal, diff.StatusChanged),
	}
	_ = tr.Apply(entries)
	if entries[0].OldVal != origVal {
		t.Error("original entry was mutated")
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	tr := truncate.New(truncate.DefaultConfig())
	out := tr.Apply([]diff.Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(out))
	}
}

func TestDefaultConfig_Values(t *testing.T) {
	cfg := truncate.DefaultConfig()
	if cfg.MaxLen != truncate.DefaultMaxLen {
		t.Errorf("expected MaxLen %d, got %d", truncate.DefaultMaxLen, cfg.MaxLen)
	}
	if cfg.Suffix != "..." {
		t.Errorf("expected suffix '...', got %q", cfg.Suffix)
	}
}

func TestNew_ZeroMaxLenUsesDefault(t *testing.T) {
	tr := truncate.New(truncate.Config{MaxLen: 0, Suffix: ">"})
	long := strings.Repeat("b", 200)
	out := tr.Apply([]diff.Entry{makeEntry("K", long, "", diff.StatusChanged)})
	if len(out[0].OldVal) != truncate.DefaultMaxLen+1 {
		t.Errorf("expected len %d, got %d", truncate.DefaultMaxLen+1, len(out[0].OldVal))
	}
}
