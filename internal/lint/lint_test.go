package lint_test

import (
	"testing"

	"github.com/you/stackdiff/internal/diff"
	"github.com/you/stackdiff/internal/lint"
)

func entry(key, newVal string, status diff.Status) diff.Entry {
	return diff.Entry{Key: key, OldValue: "", NewValue: newVal, Status: status}
}

func TestApply_NoFindings(t *testing.T) {
	entries := []diff.Entry{
		entry("HOST", "localhost", diff.StatusEqual),
		entry("PORT", "8080", diff.StatusEqual),
	}
	findings := lint.Apply(entries, lint.DefaultRules())
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestRulePlaceholder_Triggered(t *testing.T) {
	e := entry("API_URL", "<your-api-url>", diff.StatusChanged)
	f := lint.RulePlaceholder(e)
	if f == nil {
		t.Fatal("expected finding for placeholder value")
	}
	if f.Rule != "placeholder" {
		t.Errorf("unexpected rule: %s", f.Rule)
	}
	if f.Severity != lint.Error {
		t.Errorf("expected Error severity, got %s", f.Severity)
	}
}

func TestRulePlaceholder_NotTriggered(t *testing.T) {
	e := entry("API_URL", "https://example.com", diff.StatusEqual)
	if f := lint.RulePlaceholder(e); f != nil {
		t.Errorf("unexpected finding: %v", f)
	}
}

func TestRuleEmptyValue_Triggered(t *testing.T) {
	e := entry("DB_PASS", "", diff.StatusChanged)
	f := lint.RuleEmptyValue(e)
	if f == nil {
		t.Fatal("expected finding for empty value")
	}
	if f.Severity != lint.Warn {
		t.Errorf("expected Warn severity, got %s", f.Severity)
	}
}

func TestRuleEmptyValue_SkipsRemoved(t *testing.T) {
	e := entry("OLD_KEY", "", diff.StatusRemoved)
	if f := lint.RuleEmptyValue(e); f != nil {
		t.Errorf("removed entry should not trigger empty-value rule")
	}
}

func TestRuleShortSecret_Triggered(t *testing.T) {
	e := entry("API_TOKEN", "abc", diff.StatusChanged)
	f := lint.RuleShortSecret(e)
	if f == nil {
		t.Fatal("expected finding for short secret")
	}
	if f.Rule != "short-secret" {
		t.Errorf("unexpected rule: %s", f.Rule)
	}
}

func TestRuleShortSecret_LongEnough(t *testing.T) {
	e := entry("API_TOKEN", "supersecretvalue", diff.StatusEqual)
	if f := lint.RuleShortSecret(e); f != nil {
		t.Errorf("unexpected finding for long secret")
	}
}

func TestApply_MultipleFindings(t *testing.T) {
	entries := []diff.Entry{
		entry("DB_PASSWORD", "x", diff.StatusChanged),
		entry("HOST", "<host>", diff.StatusAdded),
	}
	findings := lint.Apply(entries, lint.DefaultRules())
	if len(findings) < 2 {
		t.Errorf("expected at least 2 findings, got %d", len(findings))
	}
}

func TestFinding_String(t *testing.T) {
	f := lint.Finding{Key: "K", Value: "v", Rule: "test-rule", Severity: lint.Warn, Message: "test msg"}
	s := f.String()
	if s == "" {
		t.Error("String() returned empty")
	}
}
