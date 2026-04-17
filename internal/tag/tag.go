// Package tag provides utilities for tagging diff entries with metadata labels.
package tag

import (
	"strings"

	"github.com/you/stackdiff/internal/diff"
)

// Rule maps a key prefix or pattern to a tag label.
type Rule struct {
	Prefix string
	Tag    string
}

// Tagger applies tag rules to diff entries.
type Tagger struct {
	rules []Rule
}

// New creates a Tagger with the given rules.
func New(rules []Rule) *Tagger {
	return &Tagger{rules: rules}
}

// Apply returns a map of entry key to matched tags.
func (t *Tagger) Apply(entries []diff.Entry) map[string][]string {
	result := make(map[string][]string)
	for _, e := range entries {
		var tags []string
		for _, r := range t.rules {
			if strings.HasPrefix(strings.ToLower(e.Key), strings.ToLower(r.Prefix)) {
				tags = append(tags, r.Tag)
			}
		}
		if len(tags) > 0 {
			result[e.Key] = tags
		}
	}
	return result
}

// DefaultRules returns a sensible set of built-in tagging rules.
func DefaultRules() []Rule {
	return []Rule{
		{Prefix: "db_", Tag: "database"},
		{Prefix: "redis_", Tag: "cache"},
		{Prefix: "aws_", Tag: "cloud"},
		{Prefix: "log_", Tag: "logging"},
		{Prefix: "feature_", Tag: "feature-flag"},
	}
}
