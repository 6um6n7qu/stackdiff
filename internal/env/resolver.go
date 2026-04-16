package env

import (
	"fmt"
	"os"
	"strings"
)

// ResolveResult holds the outcome of resolving a single key.
type ResolveResult struct {
	Key      string
	Value    string
	Resolved bool // true if value was substituted from environment
}

// Resolve walks a map of config entries and substitutes values that match
// the pattern ${ENV_VAR} with the corresponding OS environment variable.
// If a referenced variable is not set, an error is returned.
func Resolve(entries map[string]string) (map[string]string, []ResolveResult, error) {
	resolved := make(map[string]string, len(entries))
	results := make([]ResolveResult, 0, len(entries))

	for k, v := range entries {
		expanded, didResolve, err := expandValue(v)
		if err != nil {
			return nil, nil, fmt.Errorf("key %q: %w", k, err)
		}
		resolved[k] = expanded
		results = append(results, ResolveResult{
			Key:      k,
			Value:    expanded,
			Resolved: didResolve,
		})
	}

	return resolved, results, nil
}

// expandValue replaces ${VAR} placeholders in s with OS env values.
// Returns the expanded string, whether any substitution occurred, and any error.
func expandValue(s string) (string, bool, error) {
	if !strings.Contains(s, "${") {
		return s, false, nil
	}

	var didResolve bool
	var expandErr error

	expanded := os.Expand(s, func(key string) string {
		if expandErr != nil {
			return ""
		}
		// Only handle ${VAR} style (os.Expand strips the braces for us)
		val, ok := os.LookupEnv(key)
		if !ok {
			expandErr = fmt.Errorf("environment variable %q is not set", key)
			return ""
		}
		didResolve = true
		return val
	})

	if expandErr != nil {
		return "", false, expandErr
	}

	return expanded, didResolve, nil
}
