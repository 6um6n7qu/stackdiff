package redact

import (
	"testing"

	"github.com/user/stackdiff/internal/diff"
)

func pipelineEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "DB_PASSWORD", OldVal: "secret", NewVal: "newsecret", Status: diff.StatusChanged},
		{Key: "APP_HOST", OldVal: "localhost", NewVal: "prod.example.com", Status: diff.StatusChanged},
		{Key: "API_TOKEN", OldVal: "", NewVal: "tok123", Status: diff.StatusAdded},
		{Key: "LOG_LEVEL", OldVal: "debug", NewVal: "", Status: diff.StatusRemoved},
	}
}

func TestApplyToEntries_MasksSensitive(t *testing.T) {
	r := New()
	entries := pipelineEntries()
	result := ApplyToEntries(r, entries)

	if result[0].OldVal == "secret" {
		t.Error("expected DB_PASSWORD OldVal to be redacted")
	}
	if result[0].NewVal == "newsecret" {
		t.Error("expected DB_PASSWORD NewVal to be redacted")
	}
	if result[2].NewVal == "tok123" {
		t.Error("expected API_TOKEN NewVal to be redacted")
	}
}

func TestApplyToEntries_PreservesNonSensitive(t *testing.T) {
	r := New()
	entries := pipelineEntries()
	result := ApplyToEntries(r, entries)

	if result[1].OldVal != "localhost" {
		t.Errorf("expected APP_HOST OldVal to be unchanged, got %q", result[1].OldVal)
	}
	if result[1].NewVal != "prod.example.com" {
		t.Errorf("expected APP_HOST NewVal to be unchanged, got %q", result[1].NewVal)
	}
}

func TestApplyToEntries_EmptyValPreserved(t *testing.T) {
	r := New()
	entries := pipelineEntries()
	result := ApplyToEntries(r, entries)

	if result[2].OldVal != "" {
		t.Errorf("expected empty OldVal to remain empty, got %q", result[2].OldVal)
	}
}

func TestApplyToEntries_DoesNotMutateOriginal(t *testing.T) {
	r := New()
	entries := pipelineEntries()
	ApplyToEntries(r, entries)

	if entries[0].OldVal != "secret" {
		t.Error("original entries should not be mutated")
	}
}
