// Package classify categorises drift entries by severity based on key patterns and value changes.
package classify

import (
	"strings"

	"github.com/yourusername/stackdiff/internal/diff"
)

// Severity represents the importance level of a drift entry.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

// Result holds a drift entry alongside its assigned severity.
type Result struct {
	Entry    diff.Entry
	Severity Severity
}

// Classifier assigns severity levels to drift entries.
type Classifier struct {
	rules []Rule
}

// Rule maps a matcher function to a severity.
type Rule struct {
	Match    func(e diff.Entry) bool
	Severity Severity
}

// DefaultRules returns a sensible default rule set.
func DefaultRules() []Rule {
	return []Rule{
		{
			Match: func(e diff.Entry) bool {
				k := strings.ToLower(e.Key)
				return strings.Contains(k, "secret") || strings.Contains(k, "password") || strings.Contains(k, "token")
			},
			Severity: SeverityCritical,
		},
		{
			Match: func(e diff.Entry) bool {
				return e.Status == diff.StatusRemoved
			},
			Severity: SeverityHigh,
		},
		{
			Match: func(e diff.Entry) bool {
				return e.Status == diff.StatusAdded
			},
			Severity: SeverityMedium,
		},
		{
			Match: func(e diff.Entry) bool {
				return e.Status == diff.StatusChanged
			},
			Severity: SeverityLow,
		},
	}
}

// New creates a Classifier with the provided rules.
func New(rules []Rule) *Classifier {
	return &Classifier{rules: rules}
}

// Apply classifies each entry and returns Results. Entries with no matching rule get SeverityLow.
func (c *Classifier) Apply(entries []diff.Entry) []Result {
	out := make([]Result, 0, len(entries))
	for _, e := range entries {
		sev := SeverityLow
		for _, r := range c.rules {
			if r.Match(e) {
				sev = r.Severity
				break
			}
		}
		out = append(out, Result{Entry: e, Severity: sev})
	}
	return out
}
