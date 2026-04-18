// Package annotate attaches metadata labels to diff entries.
package annotate

import "github.com/stackdiff/stackdiff/internal/diff"

// Annotation holds key/value metadata for a diff entry.
type Annotation struct {
	Key   string
	Label string
	Meta  map[string]string
}

// Annotator applies annotations to entries based on rules.
type Annotator struct {
	rules []Rule
}

// Rule maps a condition to a set of metadata labels.
type Rule struct {
	Match  func(e diff.Entry) bool
	Labels map[string]string
}

// New creates an Annotator with the given rules.
func New(rules []Rule) *Annotator {
	return &Annotator{rules: rules}
}

// Apply returns annotations for all matching entries.
func (a *Annotator) Apply(entries []diff.Entry) []Annotation {
	var out []Annotation
	for _, e := range entries {
		for _, r := range a.rules {
			if r.Match(e) {
				out = append(out, Annotation{
					Key:   e.Key,
					Label: labelString(r.Labels),
					Meta:  r.Labels,
				})
				break
			}
		}
	}
	return out
}

func labelString(m map[string]string) string {
	for _, v := range m {
		return v
	}
	return ""
}
