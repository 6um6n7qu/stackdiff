// Package normalize provides value normalization for config entries
// before comparison, such as trimming whitespace, lowercasing booleans,
// and canonicalizing common value formats.
package normalize

import (
	"strings"

	"github.com/user/stackdiff/internal/diff"
)

// Config controls which normalizations are applied.
type Config struct {
	TrimSpace       bool
	LowercaseBools  bool
	CanonicalizeURL bool
}

// DefaultConfig returns a Config with sensible defaults enabled.
func DefaultConfig() Config {
	return Config{
		TrimSpace:       true,
		LowercaseBools:  true,
		CanonicalizeURL: false,
	}
}

// Normalizer applies normalization rules to entry values.
type Normalizer struct {
	cfg Config
}

// New creates a Normalizer with the given config.
func New(cfg Config) *Normalizer {
	return &Normalizer{cfg: cfg}
}

// Apply returns a new slice of entries with values normalized.
func (n *Normalizer) Apply(entries []diff.Entry) []diff.Entry {
	out := make([]diff.Entry, len(entries))
	for i, e := range entries {
		out[i] = diff.Entry{
			Key:    e.Key,
			OldVal: n.normalize(e.OldVal),
			NewVal: n.normalize(e.NewVal),
			Status: e.Status,
		}
	}
	return out
}

func (n *Normalizer) normalize(v string) string {
	if n.cfg.TrimSpace {
		v = strings.TrimSpace(v)
	}
	if n.cfg.LowercaseBools {
		switch strings.ToLower(v) {
		case "true", "false", "yes", "no":
			v = strings.ToLower(v)
		}
	}
	if n.cfg.CanonicalizeURL {
		v = strings.TrimRight(v, "/")
	}
	return v
}
