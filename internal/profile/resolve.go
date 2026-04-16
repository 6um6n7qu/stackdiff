package profile

import (
	"fmt"
	"os"
	"strings"
)

// ResolveResult holds the resolved profile name and its source.
type ResolveResult struct {
	Name   string
	Source string // "flag", "env", "default"
}

// Resolve determines the active profile name from the provided flag value,
// the STACKDIFF_PROFILE environment variable, or falls back to "default".
func Resolve(flagValue string) ResolveResult {
	if flagValue != "" {
		return ResolveResult{Name: normalize(flagValue), Source: "flag"}
	}
	if env := os.Getenv("STACKDIFF_PROFILE"); env != "" {
		return ResolveResult{Name: normalize(env), Source: "env"}
	}
	return ResolveResult{Name: "default", Source: "default"}
}

// ResolveFromStore resolves the active profile and loads it from the store.
// Returns an error if the profile does not exist.
func ResolveFromStore(store *Store, flagValue string) (*Profile, ResolveResult, error) {
	res := Resolve(flagValue)
	p, err := store.Load(res.Name)
	if err != nil {
		return nil, res, fmt.Errorf("profile %q not found (source: %s): %w", res.Name, res.Source, err)
	}
	return p, res, nil
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
