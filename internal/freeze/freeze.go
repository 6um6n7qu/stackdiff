// Package freeze provides utilities for detecting and enforcing frozen
// configuration keys — keys whose values must not change across environments.
package freeze

import (
	"fmt"

	"github.com/user/stackdiff/internal/diff"
)

// Violation represents a frozen key whose value has drifted.
type Violation struct {
	Key      string
	Expected string
	Actual   string
}

// String returns a human-readable description of the violation.
func (v Violation) String() string {
	return fmt.Sprintf("frozen key %q: expected %q, got %q", v.Key, v.Expected, v.Actual)
}

// Config holds the set of keys that must remain frozen.
type Config struct {
	// Keys is the list of key names that must not change.
	Keys []string
}

// DefaultConfig returns a Config with no frozen keys.
func DefaultConfig() Config {
	return Config{}
}

// Enforce checks the provided entries against the frozen key list and returns
// any violations where a frozen key has been added, removed, or changed.
func Enforce(cfg Config, entries []diff.Entry) []Violation {
	frozen := make(map[string]struct{}, len(cfg.Keys))
	for _, k := range cfg.Keys {
		frozen[k] = struct{}{}
	}

	var violations []Violation
	for _, e := range entries {
		if _, ok := frozen[e.Key]; !ok {
			continue
		}
		if e.IsDrift() {
			violations = append(violations, Violation{
				Key:      e.Key,
				Expected: e.OldValue,
				Actual:   e.NewValue,
			})
		}
	}
	return violations
}

// HasViolations returns true when at least one violation is present.
func HasViolations(vs []Violation) bool {
	return len(vs) > 0
}
