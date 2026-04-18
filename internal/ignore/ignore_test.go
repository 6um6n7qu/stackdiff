package ignore_test

import (
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/ignore"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "APP_ENV", OldVal: "staging", NewVal: "production", Status: diff.StatusChanged},
		{Key: "DEBUG_VERBOSE", OldVal: "", NewVal: "true", Status: diff.StatusAdded},
		{Key: "SECRET_KEY", OldVal: "abc", NewVal: "", Status: diff.StatusRemoved},
		{Key: "LOG_LEVEL", OldVal: "info", NewVal: "info", Status: diff.StatusEqual},
	}
}

func TestMatch_ExactKey(t *testing.T) {
	l := ignore.New([]string{"APP_ENV"})
	if !l.Match("APP_ENV") {
		t.Error("expected APP_ENV to match")
	}
	if l.Match("LOG_LEVEL") {
		t.Error("expected LOG_LEVEL not to match")
	}
}

func TestMatch_PrefixWildcard(t *testing.T) {
	l := ignore.New([]string{"DEBUG_*"})
	if !l.Match("DEBUG_VERBOSE") {
		t.Error("expected DEBUG_VERBOSE to match prefix DEBUG_*")
	}
	if l.Match("APP_ENV") {
		t.Error("expected APP_ENV not to match prefix DEBUG_*")
	}
}

func TestApply_RemovesMatchedKeys(t *testing.T) {
	l := ignore.New([]string{"APP_ENV", "DEBUG_*"})
	result := l.Apply(sampleEntries())
	for _, e := range result {
		if e.Key == "APP_ENV" || e.Key == "DEBUG_VERBOSE" {
			t.Errorf("key %q should have been ignored", e.Key)
		}
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestApply_EmptyList(t *testing.T) {
	l := ignore.New(nil)
	result := l.Apply(sampleEntries())
	if len(result) != len(sampleEntries()) {
		t.Errorf("expected all entries preserved, got %d", len(result))
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	original := sampleEntries()
	l := ignore.New([]string{"APP_ENV"})
	_ = l.Apply(original)
	if len(original) != 4 {
		t.Error("original slice was mutated")
	}
}

func TestMatch_NoPatterns(t *testing.T) {
	l := ignore.New([]string{})
	if l.Match("ANYTHING") {
		t.Error("empty ignore list should not match anything")
	}
}
