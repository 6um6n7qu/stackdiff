package annotate

import (
	"strings"

	"github.com/stackdiff/stackdiff/internal/diff"
)

// DefaultRules returns a set of built-in annotation rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Match: func(e diff.Entry) bool {
				return e.Status == diff.StatusAdded
			},
			Labels: map[string]string{"status": "added", "severity": "info"},
		},
		{
			Match: func(e diff.Entry) bool {
				return e.Status == diff.StatusRemoved
			},
			Labels: map[string]string{"status": "removed", "severity": "warning"},
		},
		{
			Match: func(e diff.Entry) bool {
				return e.Status == diff.StatusChanged
			},
			Labels: map[string]string{"status": "changed", "severity": "warning"},
		},
		{
			Match: func(e diff.Entry) bool {
				k := strings.ToLower(e.Key)
				return strings.Contains(k, "secret") || strings.Contains(k, "password") || strings.Contains(k, "token")
			},
			Labels: map[string]string{"category": "sensitive", "severity": "critical"},
		},
	}
}
