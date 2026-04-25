// Package suppress provides a mechanism to silence known or expected drift
// entries based on configurable rules, preventing noisy alerts for
// intentional configuration differences.
package suppress

import (
	"strings"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// Rule defines a suppression rule that matches drift entries.
type Rule struct {
	// Key is the exact key or prefix pattern (use "*" suffix for prefix match).
	Key string
	// Reason is a human-readable explanation for the suppression.
	Reason string
	// ExpiresAt, if non-zero, limits how long the rule is active.
	ExpiresAt time.Time
}

// Suppressor filters out drift entries matched by registered rules.
type Suppressor struct {
	rules []Rule
	now   func() time.Time
}

// New returns a Suppressor with the given rules.
func New(rules []Rule) *Suppressor {
	return &Suppressor{
		rules: rules,
		now:   time.Now,
	}
}

// Apply returns entries that are NOT suppressed by any active rule.
func (s *Suppressor) Apply(entries []diff.Entry) []diff.Entry {
	result := make([]diff.Entry, 0, len(entries))
	for _, e := range entries {
		if !s.isSuppressed(e) {
			result = append(result, e)
		}
	}
	return result
}

// CountSuppressed returns the number of entries that would be suppressed.
func (s *Suppressor) CountSuppressed(entries []diff.Entry) int {
	count := 0
	for _, e := range entries {
		if s.isSuppressed(e) {
			count++
		}
	}
	return count
}

func (s *Suppressor) isSuppressed(e diff.Entry) bool {
	now := s.now()
	for _, r := range s.rules {
		if !r.ExpiresAt.IsZero() && now.After(r.ExpiresAt) {
			continue
		}
		if matchesKey(r.Key, e.Key) {
			return true
		}
	}
	return false
}

func matchesKey(pattern, key string) bool {
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
	}
	return pattern == key
}
