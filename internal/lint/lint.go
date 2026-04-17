// Package lint provides rule-based validation for config entries,
// flagging suspicious values such as placeholders or overly short secrets.
package lint

import (
	"fmt"
	"strings"

	"github.com/you/stackdiff/internal/diff"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	Warn  Severity = "warn"
	Error Severity = "error"
)

// Finding represents a single lint result for one config entry.
type Finding struct {
	Key      string
	Value    string
	Rule     string
	Severity Severity
	Message  string
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s (%s)", f.Severity, f.Key, f.Message, f.Rule)
}

// Rule is a function that inspects an entry and returns a Finding if triggered.
type Rule func(e diff.Entry) *Finding

// DefaultRules returns the built-in set of lint rules.
func DefaultRules() []Rule {
	return []Rule{
		RulePlaceholder,
		RuleEmptyValue,
		RuleShortSecret,
	}
}

// RulePlaceholder flags values that look like unfilled template placeholders.
func RulePlaceholder(e diff.Entry) *Finding {
	v := e.NewValue
	if strings.HasPrefix(v, "<") && strings.HasSuffix(v, ">") {
		return &Finding{Key: e.Key, Value: v, Rule: "placeholder", Severity: Error,
			Message: "value appears to be an unfilled placeholder"}
	}
	return nil
}

// RuleEmptyValue flags entries whose new value is empty.
func RuleEmptyValue(e diff.Entry) *Finding {
	if e.NewValue == "" && e.Status != diff.StatusRemoved {
		return &Finding{Key: e.Key, Value: e.NewValue, Rule: "empty-value", Severity: Warn,
			Message: "value is empty"}
	}
	return nil
}

// RuleShortSecret flags secret/password keys whose value is suspiciously short.
func RuleShortSecret(e diff.Entry) *Finding {
	k := strings.ToLower(e.Key)
	if (strings.Contains(k, "secret") || strings.Contains(k, "password") || strings.Contains(k, "token")) &&
		len(e.NewValue) > 0 && len(e.NewValue) < 8 {
		return &Finding{Key: e.Key, Value: e.NewValue, Rule: "short-secret", Severity: Warn,
			Message: "secret value is suspiciously short (< 8 chars)"}
	}
	return nil
}

// Apply runs all provided rules against each entry and returns all findings.
func Apply(entries []diff.Entry, rules []Rule) []Finding {
	var findings []Finding
	for _, e := range entries {
		for _, r := range rules {
			if f := r(e); f != nil {
				findings = append(findings, *f)
			}
		}
	}
	return findings
}
