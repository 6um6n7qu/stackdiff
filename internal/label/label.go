// Package label provides functionality for attaching and querying
// string labels (key=value metadata) on diff entries.
package label

import (
	"fmt"
	"strings"

	"github.com/yourorg/stackdiff/internal/diff"
)

// Labeler attaches metadata labels to diff entries based on rules.
type Labeler struct {
	rules []Rule
}

// Rule defines a labeling rule: if Match returns true for an entry,
// the returned labels are merged into the entry's Meta map.
type Rule struct {
	Name  string
	Match func(e diff.Entry) bool
	Label func(e diff.Entry) map[string]string
}

// New returns a Labeler configured with the provided rules.
func New(rules []Rule) *Labeler {
	return &Labeler{rules: rules}
}

// Apply attaches labels to each entry by evaluating all rules.
// It does not mutate the original slice; it returns new Entry copies.
func (l *Labeler) Apply(entries []diff.Entry) []diff.Entry {
	out := make([]diff.Entry, len(entries))
	for i, e := range entries {
		merged := cloneMeta(e.Meta)
		for _, r := range l.rules {
			if r.Match(e) {
				for k, v := range r.Label(e) {
					merged[k] = v
				}
			}
		}
		e.Meta = merged
		out[i] = e
	}
	return out
}

// DefaultRules returns a baseline set of labeling rules covering
// drift status and key sensitivity.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name:  "status",
			Match: func(e diff.Entry) bool { return true },
			Label: func(e diff.Entry) map[string]string {
				return map[string]string{"status": string(e.Status)}
			},
		},
		{
			Name: "sensitive",
			Match: func(e diff.Entry) bool {
				key := strings.ToLower(e.Key)
				for _, kw := range []string{"password", "secret", "token", "key", "apikey"} {
					if strings.Contains(key, kw) {
						return true
					}
				}
				return false
			},
			Label: func(e diff.Entry) map[string]string {
				return map[string]string{"sensitive": "true"}
			},
		},
	}
}

// Format returns a human-readable representation of an entry's labels.
func Format(meta map[string]string) string {
	if len(meta) == 0 {
		return ""
	}
	parts := make([]string, 0, len(meta))
	for k, v := range meta {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, " ")
}

func cloneMeta(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
