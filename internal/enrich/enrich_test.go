package enrich_test

import (
	"testing"

	"github.com/yourusername/stackdiff/internal/diff"
	"github.com/yourusername/stackdiff/internal/enrich"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "DB_PASSWORD", OldVal: "old", NewVal: "new", Status: diff.StatusChanged},
		{Key: "APP_PORT", OldVal: "8080", NewVal: "9090", Status: diff.StatusChanged},
		{Key: "API_TOKEN", OldVal: "", NewVal: "abc", Status: diff.StatusAdded},
	}
}

func TestApply_SourceRule(t *testing.T) {
	en := enrich.New(enrich.SourceRule("prod"))
	out := en.Apply(sampleEntries())
	for _, e := range out {
		if e.Metadata["source"] != "prod" {
			t.Errorf("expected source=prod for key %s, got %s", e.Key, e.Metadata["source"])
		}
	}
}

func TestApply_SensitiveRule_Detected(t *testing.T) {
	en := enrich.New(enrich.SensitiveRule([]string{"password", "token"}))
	out := en.Apply(sampleEntries())
	sensitiveKeys := map[string]bool{"DB_PASSWORD": true, "API_TOKEN": true}
	for _, e := range out {
		expected := "false"
		if sensitiveKeys[e.Key] {
			expected = "true"
		}
		if e.Metadata["sensitive"] != expected {
			t.Errorf("key %s: expected sensitive=%s, got %s", e.Key, expected, e.Metadata["sensitive"])
		}
	}
}

func TestApply_MultipleRules_MergesMeta(t *testing.T) {
	en := enrich.New(enrich.SourceRule("staging"), enrich.SensitiveRule([]string{"password"}))
	out := en.Apply(sampleEntries())
	for _, e := range out {
		if e.Metadata["source"] != "staging" {
			t.Errorf("missing source for %s", e.Key)
		}
		if _, ok := e.Metadata["sensitive"]; !ok {
			t.Errorf("missing sensitive for %s", e.Key)
		}
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	original := sampleEntries()
	en := enrich.New(enrich.SourceRule("test"))
	en.Apply(original)
	for _, e := range original {
		if e.Metadata != nil {
			t.Errorf("original entry %s was mutated", e.Key)
		}
	}
}

func TestDefaultRules_NotEmpty(t *testing.T) {
	rules := enrich.DefaultRules()
	if len(rules) == 0 {
		t.Error("expected at least one default rule")
	}
}
