package segment_test

import (
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/segment"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "HOST", OldValue: "a", NewValue: "b", Status: diff.StatusChanged},
		{Key: "PORT", OldValue: "", NewValue: "8080", Status: diff.StatusAdded},
		{Key: "DEBUG", OldValue: "true", NewValue: "", Status: diff.StatusRemoved},
		{Key: "APP", OldValue: "myapp", NewValue: "myapp", Status: diff.StatusEqual},
	}
}

func TestApply_DefaultRules_PlacesAllEntries(t *testing.T) {
	entries := sampleEntries()
	result, err := segment.Apply(entries, segment.DefaultRules())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Unmatched) != 0 {
		t.Errorf("expected no unmatched entries, got %d", len(result.Unmatched))
	}
	if len(result.Buckets["changed"]) != 1 {
		t.Errorf("expected 1 changed entry, got %d", len(result.Buckets["changed"]))
	}
	if len(result.Buckets["added"]) != 1 {
		t.Errorf("expected 1 added entry, got %d", len(result.Buckets["added"]))
	}
	if len(result.Buckets["removed"]) != 1 {
		t.Errorf("expected 1 removed entry, got %d", len(result.Buckets["removed"]))
	}
	if len(result.Buckets["equal"]) != 1 {
		t.Errorf("expected 1 equal entry, got %d", len(result.Buckets["equal"]))
	}
}

func TestApply_UnmatchedEntries(t *testing.T) {
	entries := sampleEntries()
	rules := []segment.Rule{
		{Name: "changed", Predicate: func(e diff.Entry) bool { return e.Status == diff.StatusChanged }},
	}
	result, err := segment.Apply(entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Unmatched) != 3 {
		t.Errorf("expected 3 unmatched entries, got %d", len(result.Unmatched))
	}
}

func TestApply_NoRules_ReturnsError(t *testing.T) {
	_, err := segment.Apply(sampleEntries(), nil)
	if err == nil {
		t.Fatal("expected error for empty rules, got nil")
	}
}

func TestResult_HasDrift_True(t *testing.T) {
	result, _ := segment.Apply(sampleEntries(), segment.DefaultRules())
	if !result.HasDrift() {
		t.Error("expected HasDrift to be true")
	}
}

func TestResult_HasDrift_False(t *testing.T) {
	entries := []diff.Entry{
		{Key: "A", OldValue: "x", NewValue: "x", Status: diff.StatusEqual},
	}
	result, _ := segment.Apply(entries, segment.DefaultRules())
	if result.HasDrift() {
		t.Error("expected HasDrift to be false")
	}
}

func TestApply_CustomRule_KeyPrefix(t *testing.T) {
	entries := []diff.Entry{
		{Key: "DB_HOST", OldValue: "a", NewValue: "b", Status: diff.StatusChanged},
		{Key: "APP_ENV", OldValue: "prod", NewValue: "prod", Status: diff.StatusEqual},
	}
	rules := []segment.Rule{
		{
			Name: "db",
			Predicate: func(e diff.Entry) bool {
				return len(e.Key) >= 3 && e.Key[:3] == "DB_"
			},
		},
	}
	result, err := segment.Apply(entries, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Buckets["db"]) != 1 {
		t.Errorf("expected 1 db entry, got %d", len(result.Buckets["db"]))
	}
	if len(result.Unmatched) != 1 {
		t.Errorf("expected 1 unmatched entry, got %d", len(result.Unmatched))
	}
}
