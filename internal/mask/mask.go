package mask

import (
	"strings"
)

// Rule defines a masking rule for a specific key pattern.
type Rule struct {
	Pattern     string
	Replacement string
}

// Masker applies masking rules to config entry values.
type Masker struct {
	rules []Rule
}

var defaultPatterns = []string{
	"password", "secret", "token", "key", "apikey", "api_key",
	"auth", "credential", "private", "passphrase",
}

const DefaultReplacement = "***"

// New creates a Masker with default sensitive key patterns.
func New() *Masker {
	rules := make([]Rule, 0, len(defaultPatterns))
	for _, p := range defaultPatterns {
		rules = append(rules, Rule{Pattern: p, Replacement: DefaultReplacement})
	}
	return &Masker{rules: rules}
}

// NewWithRules creates a Masker with custom rules.
func NewWithRules(rules []Rule) *Masker {
	return &Masker{rules: rules}
}

// AddRule appends a masking rule.
func (m *Masker) AddRule(r Rule) {
	m.rules = append(m.rules, r)
}

// IsSensitive returns true if the key matches any masking rule.
func (m *Masker) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, r := range m.rules {
		if strings.Contains(lower, strings.ToLower(r.Pattern)) {
			return true
		}
	}
	return false
}

// Apply masks values in the provided map for sensitive keys.
func (m *Masker) Apply(entries map[string]string) map[string]string {
	result := make(map[string]string, len(entries))
	for k, v := range entries {
		if m.IsSensitive(k) {
			result[k] = m.replacementFor(k)
		} else {
			result[k] = v
		}
	}
	return result
}

func (m *Masker) replacementFor(key string) string {
	lower := strings.ToLower(key)
	for _, r := range m.rules {
		if strings.Contains(lower, strings.ToLower(r.Pattern)) {
			return r.Replacement
		}
	}
	return DefaultReplacement
}
