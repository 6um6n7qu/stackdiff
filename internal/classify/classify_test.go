package classify_test

import (
	"testing"

	"github.com/yourusername/stackdiff/internal/classify"
	"github.com/yourusername/stackdiff/internal/diff"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "DB_PASSWORD", OldVal: "old", NewVal: "new", Status: diff.StatusChanged},
		{Key: "APP_ENV", OldVal: "", NewVal: "prod", Status: diff.StatusAdded},
		{Key: "LEGACY_KEY", OldVal: "val", NewVal: "", Status: diff.StatusRemoved},
		{Key: "LOG_LEVEL", OldVal: "info", NewVal: "debug", Status: diff.StatusChanged},
	}
}

func TestApply_CriticalForSecretKey(t *testing.T) {
	c := classify.New(classify.DefaultRules())
	results := c.Apply(sampleEntries())
	if results[0].Severity != classify.SeverityCritical {
		t.Errorf("expected critical for password key, got %s", results[0].Severity)
	}
}

func TestApply_HighForRemoved(t *testing.T) {
	c := classify.New(classify.DefaultRules())
	results := c.Apply(sampleEntries())
	if results[2].Severity != classify.SeverityHigh {
		t.Errorf("expected high for removed entry, got %s", results[2].Severity)
	}
}

func TestApply_MediumForAdded(t *testing.T) {
	c := classify.New(classify.DefaultRules())
	results := c.Apply(sampleEntries())
	if results[1].Severity != classify.SeverityMedium {
		t.Errorf("expected medium for added entry, got %s", results[1].Severity)
	}
}

func TestApply_LowForChangedNonSensitive(t *testing.T) {
	c := classify.New(classify.DefaultRules())
	results := c.Apply(sampleEntries())
	if results[3].Severity != classify.SeverityLow {
		t.Errorf("expected low for changed non-sensitive entry, got %s", results[3].Severity)
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	c := classify.New(classify.DefaultRules())
	results := c.Apply([]diff.Entry{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestApply_CustomRule(t *testing.T) {
	rules := []classify.Rule{
		{
			Match:    func(e diff.Entry) bool { return e.Key == "SPECIAL" },
			Severity: classify.SeverityCritical,
		},
	}
	c := classify.New(rules)
	entries := []diff.Entry{
		{Key: "SPECIAL", OldVal: "a", NewVal: "b", Status: diff.StatusChanged},
		{Key: "OTHER", OldVal: "a", NewVal: "b", Status: diff.StatusChanged},
	}
	results := c.Apply(entries)
	if results[0].Severity != classify.SeverityCritical {
		t.Errorf("expected critical for SPECIAL, got %s", results[0].Severity)
	}
	if results[1].Severity != classify.SeverityLow {
		t.Errorf("expected low fallback for OTHER, got %s", results[1].Severity)
	}
}
