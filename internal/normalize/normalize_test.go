package normalize_test

import (
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/normalize"
)

func makeEntry(key, oldVal, newVal string) diff.Entry {
	return diff.Entry{Key: key, OldVal: oldVal, NewVal: newVal, Status: diff.StatusChanged}
}

func TestApply_TrimSpace(t *testing.T) {
	n := normalize.New(normalize.DefaultConfig())
	entries := []diff.Entry{makeEntry("key", "  hello  ", " world ")}
	out := n.Apply(entries)
	if out[0].OldVal != "hello" {
		t.Errorf("expected 'hello', got %q", out[0].OldVal)
	}
	if out[0].NewVal != "world" {
		t.Errorf("expected 'world', got %q", out[0].NewVal)
	}
}

func TestApply_LowercaseBools(t *testing.T) {
	n := normalize.New(normalize.DefaultConfig())
	entries := []diff.Entry{makeEntry("enabled", "True", "FALSE")}
	out := n.Apply(entries)
	if out[0].OldVal != "true" {
		t.Errorf("expected 'true', got %q", out[0].OldVal)
	}
	if out[0].NewVal != "false" {
		t.Errorf("expected 'false', got %q", out[0].NewVal)
	}
}

func TestApply_CanonicalizeURL(t *testing.T) {
	cfg := normalize.DefaultConfig()
	cfg.CanonicalizeURL = true
	n := normalize.New(cfg)
	entries := []diff.Entry{makeEntry("url", "http://example.com/", "http://other.com//")}
	out := n.Apply(entries)
	if out[0].OldVal != "http://example.com" {
		t.Errorf("unexpected OldVal: %q", out[0].OldVal)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	n := normalize.New(normalize.DefaultConfig())
	orig := []diff.Entry{makeEntry("k", "  v  ", "  w  ")}
	n.Apply(orig)
	if orig[0].OldVal != "  v  " {
		t.Error("original entry was mutated")
	}
}

func TestApplyDefault_Pipeline(t *testing.T) {
	entries := []diff.Entry{makeEntry("x", " YES ", " NO ")}
	out := normalize.ApplyDefault(entries)
	if out[0].OldVal != "yes" {
		t.Errorf("expected 'yes', got %q", out[0].OldVal)
	}
}

func TestDefaultConfig_Fields(t *testing.T) {
	cfg := normalize.DefaultConfig()
	if !cfg.TrimSpace {
		t.Error("expected TrimSpace to be true")
	}
	if !cfg.LowercaseBools {
		t.Error("expected LowercaseBools to be true")
	}
	if cfg.CanonicalizeURL {
		t.Error("expected CanonicalizeURL to be false by default")
	}
}
