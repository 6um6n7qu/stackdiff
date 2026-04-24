package freeze_test

import (
	"strings"
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/freeze"
)

func makeEntry(key, oldVal, newVal string, status diff.Status) diff.Entry {
	return diff.Entry{Key: key, OldValue: oldVal, NewValue: newVal, Status: status}
}

func TestEnforce_NoFrozenKeys(t *testing.T) {
	cfg := freeze.DefaultConfig()
	entries := []diff.Entry{
		makeEntry("APP_ENV", "staging", "production", diff.StatusChanged),
	}
	vs := freeze.Enforce(cfg, entries)
	if len(vs) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(vs))
	}
}

func TestEnforce_FrozenKeyUnchanged(t *testing.T) {
	cfg := freeze.Config{Keys: []string{"DB_HOST"}}
	entries := []diff.Entry{
		makeEntry("DB_HOST", "db.prod", "db.prod", diff.StatusEqual),
	}
	vs := freeze.Enforce(cfg, entries)
	if len(vs) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(vs))
	}
}

func TestEnforce_FrozenKeyChanged(t *testing.T) {
	cfg := freeze.Config{Keys: []string{"DB_HOST"}}
	entries := []diff.Entry{
		makeEntry("DB_HOST", "db.prod", "db.staging", diff.StatusChanged),
	}
	vs := freeze.Enforce(cfg, entries)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Key != "DB_HOST" {
		t.Errorf("unexpected key %q", vs[0].Key)
	}
}

func TestEnforce_FrozenKeyAdded(t *testing.T) {
	cfg := freeze.Config{Keys: []string{"NEW_KEY"}}
	entries := []diff.Entry{
		makeEntry("NEW_KEY", "", "value", diff.StatusAdded),
	}
	vs := freeze.Enforce(cfg, entries)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
}

func TestEnforce_NonFrozenKeyIgnored(t *testing.T) {
	cfg := freeze.Config{Keys: []string{"DB_HOST"}}
	entries := []diff.Entry{
		makeEntry("APP_ENV", "staging", "production", diff.StatusChanged),
		makeEntry("DB_HOST", "db.prod", "db.prod", diff.StatusEqual),
	}
	vs := freeze.Enforce(cfg, entries)
	if len(vs) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(vs))
	}
}

func TestMustPass_ReturnsNilWhenOK(t *testing.T) {
	cfg := freeze.Config{Keys: []string{"DB_HOST"}}
	entries := []diff.Entry{
		makeEntry("DB_HOST", "db.prod", "db.prod", diff.StatusEqual),
	}
	if err := freeze.MustPass(cfg, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMustPass_ReturnsErrorOnBreach(t *testing.T) {
	cfg := freeze.Config{Keys: []string{"DB_HOST"}}
	entries := []diff.Entry{
		makeEntry("DB_HOST", "db.prod", "db.staging", diff.StatusChanged),
	}
	err := freeze.MustPass(cfg, entries)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "DB_HOST") {
		t.Errorf("error should mention key, got: %v", err)
	}
}

func TestResult_Summary_NoDrift(t *testing.T) {
	r := freeze.Run(freeze.DefaultConfig(), nil)
	if !strings.Contains(r.Summary(), "no frozen-key") {
		t.Errorf("unexpected summary: %q", r.Summary())
	}
}
