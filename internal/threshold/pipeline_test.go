package threshold_test

import (
	"testing"

	"github.com/yourusername/stackdiff/internal/diff"
	"github.com/yourusername/stackdiff/internal/threshold"
)

func driftEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "DB_HOST", OldValue: "localhost", NewValue: "prod-db", Status: diff.StatusChanged},
		{Key: "API_KEY", OldValue: "", NewValue: "abc123", Status: diff.StatusAdded},
		{Key: "OLD_FLAG", OldValue: "true", NewValue: "", Status: diff.StatusRemoved},
	}
}

func TestCheckEntries_Breached(t *testing.T) {
	entries := driftEntries()
	cfg := threshold.Config{Warning: 1.0, Critical: 100.0}
	r := threshold.CheckEntries(entries, cfg)
	if !r.Breached() {
		t.Error("expected threshold to be breached with low warning limit")
	}
}

func TestCheckEntries_NoDrift(t *testing.T) {
	entries := []diff.Entry{
		{Key: "HOST", OldValue: "a", NewValue: "a", Status: diff.StatusEqual},
	}
	cfg := threshold.DefaultConfig()
	r := threshold.CheckEntries(entries, cfg)
	if r.Breached() {
		t.Errorf("expected no breach for equal entries, got level %s", r.Level)
	}
}

func TestMustNotBreach_ReturnsEmptyWhenOK(t *testing.T) {
	entries := []diff.Entry{
		{Key: "X", OldValue: "1", NewValue: "1", Status: diff.StatusEqual},
	}
	msg := threshold.MustNotBreach(entries, threshold.DefaultConfig())
	if msg != "" {
		t.Errorf("expected empty message, got %q", msg)
	}
}

func TestMustNotBreach_ReturnsMessageWhenBreached(t *testing.T) {
	cfg := threshold.Config{Warning: 0.1, Critical: 200.0}
	msg := threshold.MustNotBreach(driftEntries(), cfg)
	if msg == "" {
		t.Error("expected breach message, got empty string")
	}
}
