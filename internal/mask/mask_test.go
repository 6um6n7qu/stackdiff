package mask

import (
	"testing"
)

func TestIsSensitive_MatchesPassword(t *testing.T) {
	m := New()
	if !m.IsSensitive("db_password") {
		t.Error("expected db_password to be sensitive")
	}
}

func TestIsSensitive_MatchesToken(t *testing.T) {
	m := New()
	if !m.IsSensitive("AUTH_TOKEN") {
		t.Error("expected AUTH_TOKEN to be sensitive")
	}
}

func TestIsSensitive_SafeKey(t *testing.T) {
	m := New()
	if m.IsSensitive("log_level") {
		t.Error("expected log_level to not be sensitive")
	}
}

func TestIsSensitive_CaseInsensitive(t *testing.T) {
	m := New()
	if !m.IsSensitive("API_KEY") {
		t.Error("expected API_KEY to be sensitive")
	}
}

func TestApply_MasksSensitiveValues(t *testing.T) {
	m := New()
	input := map[string]string{
		"db_password": "supersecret",
		"log_level":   "info",
	}
	out := m.Apply(input)
	if out["db_password"] != DefaultReplacement {
		t.Errorf("expected masked value, got %q", out["db_password"])
	}
	if out["log_level"] != "info" {
		t.Errorf("expected 'info', got %q", out["log_level"])
	}
}

func TestApply_PreservesNonSensitive(t *testing.T) {
	m := New()
	input := map[string]string{"port": "8080", "host": "localhost"}
	out := m.Apply(input)
	for k, v := range input {
		if out[k] != v {
			t.Errorf("key %q: expected %q got %q", k, v, out[k])
		}
	}
}

func TestNewWithRules_CustomReplacement(t *testing.T) {
	rules := []Rule{{Pattern: "secret", Replacement: "<hidden>"}}
	m := NewWithRules(rules)
	input := map[string]string{"my_secret": "val", "other": "val2"}
	out := m.Apply(input)
	if out["my_secret"] != "<hidden>" {
		t.Errorf("expected <hidden>, got %q", out["my_secret"])
	}
	if out["other"] != "val2" {
		t.Error("non-sensitive key should not be masked")
	}
}

func TestAddRule_ExtendsMasker(t *testing.T) {
	m := New()
	m.AddRule(Rule{Pattern: "internal", Replacement: "[redacted]"})
	if !m.IsSensitive("internal_flag") {
		t.Error("expected internal_flag to be sensitive after AddRule")
	}
}
