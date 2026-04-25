package suppress

import (
	"testing"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

func makeEntry(key, oldVal, newVal string, status diff.Status) diff.Entry {
	return diff.Entry{
		Key:      key,
		OldValue: oldVal,
		NewValue: newVal,
		Status:   status,
	}
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	s := New(nil)
	entries := []diff.Entry{
		makeEntry("FOO", "a", "b", diff.StatusChanged),
		makeEntry("BAR", "", "c", diff.StatusAdded),
	}
	got := s.Apply(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestApply_ExactKeyMatch(t *testing.T) {
	s := New([]Rule{{Key: "DEPLOY_SHA", Reason: "noisy"}})
	entries := []diff.Entry{
		makeEntry("DEPLOY_SHA", "old", "new", diff.StatusChanged),
		makeEntry("APP_ENV", "staging", "prod", diff.StatusChanged),
	}
	got := s.Apply(entries)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV, got %s", got[0].Key)
	}
}

func TestApply_PrefixWildcard(t *testing.T) {
	s := New([]Rule{{Key: "TMP_*", Reason: "temporary"}})
	entries := []diff.Entry{
		makeEntry("TMP_TOKEN", "a", "b", diff.StatusChanged),
		makeEntry("TMP_FLAG", "", "1", diff.StatusAdded),
		makeEntry("REAL_KEY", "x", "y", diff.StatusChanged),
	}
	got := s.Apply(entries)
	if len(got) != 1 || got[0].Key != "REAL_KEY" {
		t.Errorf("expected only REAL_KEY, got %+v", got)
	}
}

func TestApply_ExpiredRuleDoesNotSuppress(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	s := New([]Rule{{Key: "DEPLOY_SHA", Reason: "expired", ExpiresAt: past}})
	entries := []diff.Entry{
		makeEntry("DEPLOY_SHA", "old", "new", diff.StatusChanged),
	}
	got := s.Apply(entries)
	if len(got) != 1 {
		t.Errorf("expected entry to pass through expired rule, got %d entries", len(got))
	}
}

func TestApply_ActiveExpiryStillSuppresses(t *testing.T) {
	future := time.Now().Add(1 * time.Hour)
	s := New([]Rule{{Key: "DEPLOY_SHA", Reason: "soon", ExpiresAt: future}})
	entries := []diff.Entry{
		makeEntry("DEPLOY_SHA", "old", "new", diff.StatusChanged),
	}
	got := s.Apply(entries)
	if len(got) != 0 {
		t.Errorf("expected entry to be suppressed, got %d entries", len(got))
	}
}

func TestCountSuppressed(t *testing.T) {
	s := New([]Rule{
		{Key: "DEPLOY_SHA"},
		{Key: "BUILD_*"},
	})
	entries := []diff.Entry{
		makeEntry("DEPLOY_SHA", "a", "b", diff.StatusChanged),
		makeEntry("BUILD_ID", "1", "2", diff.StatusChanged),
		makeEntry("APP_ENV", "dev", "prod", diff.StatusChanged),
	}
	count := s.CountSuppressed(entries)
	if count != 2 {
		t.Errorf("expected 2 suppressed, got %d", count)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	s := New([]Rule{{Key: "FOO"}})
	original := []diff.Entry{
		makeEntry("FOO", "a", "b", diff.StatusChanged),
		makeEntry("BAR", "c", "d", diff.StatusChanged),
	}
	copy := append([]diff.Entry{}, original...)
	s.Apply(original)
	if len(original) != len(copy) {
		t.Error("Apply mutated the original slice")
	}
}
