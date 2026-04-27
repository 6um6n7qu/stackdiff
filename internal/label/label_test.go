package label_test

import (
	"testing"

	"github.com/yourorg/stackdiff/internal/diff"
	"github.com/yourorg/stackdiff/internal/label"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "APP_ENV", OldValue: "staging", NewValue: "production", Status: diff.StatusChanged},
		{Key: "DB_PASSWORD", OldValue: "", NewValue: "s3cr3t", Status: diff.StatusAdded},
		{Key: "LOG_LEVEL", OldValue: "info", NewValue: "info", Status: diff.StatusEqual},
		{Key: "API_KEY", OldValue: "abc", NewValue: "", Status: diff.StatusRemoved},
	}
}

func TestApply_DefaultRules_StatusLabel(t *testing.T) {
	l := label.New(label.DefaultRules())
	out := l.Apply(sampleEntries())

	if got := out[0].Meta["status"]; got != string(diff.StatusChanged) {
		t.Errorf("expected status=changed, got %q", got)
	}
	if got := out[2].Meta["status"]; got != string(diff.StatusEqual) {
		t.Errorf("expected status=equal, got %q", got)
	}
}

func TestApply_DefaultRules_SensitiveLabel(t *testing.T) {
	l := label.New(label.DefaultRules())
	out := l.Apply(sampleEntries())

	if got := out[1].Meta["sensitive"]; got != "true" {
		t.Errorf("DB_PASSWORD should be sensitive, got %q", got)
	}
	if got := out[3].Meta["sensitive"]; got != "true" {
		t.Errorf("API_KEY should be sensitive, got %q", got)
	}
	if _, ok := out[0].Meta["sensitive"]; ok {
		t.Error("APP_ENV should not be marked sensitive")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	entries := sampleEntries()
	l := label.New(label.DefaultRules())
	_ = l.Apply(entries)

	for _, e := range entries {
		if len(e.Meta) != 0 {
			t.Errorf("original entry %q was mutated", e.Key)
		}
	}
}

func TestApply_CustomRule(t *testing.T) {
	rules := []label.Rule{
		{
			Name:  "env-tag",
			Match: func(e diff.Entry) bool { return e.Status == diff.StatusAdded },
			Label: func(e diff.Entry) map[string]string { return map[string]string{"event": "added"} },
		},
	}
	l := label.New(rules)
	out := l.Apply(sampleEntries())

	if got := out[1].Meta["event"]; got != "added" {
		t.Errorf("expected event=added for added entry, got %q", got)
	}
	if _, ok := out[0].Meta["event"]; ok {
		t.Error("changed entry should not have event label")
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	l := label.New(label.DefaultRules())
	out := l.Apply([]diff.Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}

func TestFormat_NonEmpty(t *testing.T) {
	meta := map[string]string{"status": "changed", "sensitive": "true"}
	result := label.Format(meta)
	if result == "" {
		t.Error("expected non-empty format output")
	}
}

func TestFormat_Empty(t *testing.T) {
	result := label.Format(map[string]string{})
	if result != "" {
		t.Errorf("expected empty string for empty meta, got %q", result)
	}
}
