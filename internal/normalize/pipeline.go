package normalize

import "github.com/user/stackdiff/internal/diff"

// ApplyDefault normalizes entries using the default configuration.
func ApplyDefault(entries []diff.Entry) []diff.Entry {
	return New(DefaultConfig()).Apply(entries)
}

// ApplyWithURL normalizes entries and also canonicalizes URL values.
func ApplyWithURL(entries []diff.Entry) []diff.Entry {
	cfg := DefaultConfig()
	cfg.CanonicalizeURL = true
	return New(cfg).Apply(entries)
}
