package filter_test

import (
	"testing"

	"github.com/yourusername/stackdiff/internal/diff"
	"github.com/yourusername/stackdiff/internal/filter"
)

var sampleEntries = []diff.Entry{
	{Key: "APP_ENV", Left: "staging", Right: "production", Status: diff.StatusChanged},
	{Key: "DB_HOST", Left: "", Right: "db.prod", Status: diff.StatusAdded},
	{Key: "LEGACY_KEY", Left: "old", Right: "", Status: diff.StatusRemoved},
	{Key: "APP_PORT", Left: "8080", Right: "8080", Status: diff.StatusEqual},
}

func TestApply_NoFilters(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{})
	if len(result) != len(sampleEntries) {
		t.Errorf("expected %d entries, got %d", len(sampleEntries), len(result))
	}
}

func TestApply_OnlyChanged(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{OnlyChanged: true})
	if len(result) != 1 || result[0].Key != "APP_ENV" {
		t.Errorf("expected only APP_ENV, got %+v", result)
	}
}

func TestApply_OnlyAdded(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{OnlyAdded: true})
	if len(result) != 1 || result[0].Key != "DB_HOST" {
		t.Errorf("expected only DB_HOST, got %+v", result)
	}
}

func TestApply_OnlyRemoved(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{OnlyRemoved: true})
	if len(result) != 1 || result[0].Key != "LEGACY_KEY" {
		t.Errorf("expected only LEGACY_KEY, got %+v", result)
	}
}

func TestApply_KeyPrefix(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{KeyPrefix: "APP_"})
	if len(result) != 2 {
		t.Errorf("expected 2 entries with prefix APP_, got %d", len(result))
	}
}

func TestApply_KeyPrefixWithOnlyChanged(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{KeyPrefix: "APP_", OnlyChanged: true})
	if len(result) != 1 || result[0].Key != "APP_ENV" {
		t.Errorf("expected only APP_ENV, got %+v", result)
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	result := filter.Apply([]diff.Entry{}, filter.Options{OnlyChanged: true})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
