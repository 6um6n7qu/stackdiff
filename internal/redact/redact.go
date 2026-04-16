// Package redact provides utilities for masking sensitive config values
// before display, export, or audit logging.
package redact

import "strings"

// DefaultPatterns are key substrings that trigger redaction.
var DefaultPatterns = []string{
	"password", "secret", "token", "apikey", "api_key", "private", "credential",
}

// Redactor masks sensitive values in config maps.
type Redactor struct {
	Patterns []string
	Mask     string
}

// New returns a Redactor with default patterns and mask.
func New() *Redactor {
	return &Redactor{
		Patterns: DefaultPatterns,
		Mask:     "***REDACTED***",
	}
}

// IsSensitive returns true if the key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.Patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// Apply returns a copy of the map with sensitive values masked.
func (r *Redactor) Apply(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if r.IsSensitive(k) {
			out[k] = r.Mask
		} else {
			out[k] = v
		}
	}
	return out
}
