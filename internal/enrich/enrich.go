// Package enrich attaches metadata to diff entries based on configurable rules.
package enrich

import (
	"strings"

	"github.com/yourusername/stackdiff/internal/diff"
)

// Rule defines a function that enriches a single entry's metadata.
type Rule func(e diff.Entry) map[string]string

// Enricher applies a set of rules to entries.
type Enricher struct {
	rules []Rule
}

// New returns an Enricher with the given rules.
func New(rules ...Rule) *Enricher {
	return &Enricher{rules: rules}
}

// Apply returns a new slice of entries with Metadata populated by the rules.
func (en *Enricher) Apply(entries []diff.Entry) []diff.Entry {
	out := make([]diff.Entry, len(entries))
	for i, e := range entries {
		meta := make(map[string]string)
		for _, r := range en.rules {
			for k, v := range r(e) {
				meta[k] = v
			}
		}
		e.Metadata = meta
		out[i] = e
	}
	return out
}

// SourceRule tags each entry with a fixed source label.
func SourceRule(source string) Rule {
	return func(e diff.Entry) map[string]string {
		return map[string]string{"source": source}
	}
}

// SensitiveRule marks entries whose keys contain sensitive substrings.
func SensitiveRule(keywords []string) Rule {
	return func(e diff.Entry) map[string]string {
		key := strings.ToLower(e.Key)
		for _, kw := range keywords {
			if strings.Contains(key, strings.ToLower(kw)) {
				return map[string]string{"sensitive": "true"}
			}
		}
		return map[string]string{"sensitive": "false"}
	}
}

// DefaultRules returns a standard set of enrichment rules.
func DefaultRules() []Rule {
	return []Rule{
		SensitiveRule([]string{"password", "secret", "token", "key", "api"}),
	}
}
