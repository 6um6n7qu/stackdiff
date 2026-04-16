package profile

import (
	"fmt"

	"github.com/stackdiff/stackdiff/internal/config"
	"github.com/stackdiff/stackdiff/internal/env"
)

// ResolvedProfile holds a loaded profile and its resolved config entries.
type ResolvedProfile struct {
	Profile *Profile
	Entries map[string]string
}

// Resolve loads the env file referenced by the profile, applies env var
// expansion, and returns a ResolvedProfile ready for comparison.
func Resolve(p *Profile) (*ResolvedProfile, error) {
	if p.EnvFile == "" {
		return &ResolvedProfile{Profile: p, Entries: map[string]string{}}, nil
	}

	raw, err := env.LoadFile(p.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("profile resolve: load env file: %w", err)
	}

	resolved, err := env.Resolve(raw)
	if err != nil {
		return nil, fmt.Errorf("profile resolve: expand vars: %w", err)
	}

	cfg := &config.Config{Entries: resolved}
	if err := config.Validate(cfg); err != nil {
		return nil, fmt.Errorf("profile resolve: validate: %w", err)
	}

	return &ResolvedProfile{Profile: p, Entries: resolved}, nil
}
