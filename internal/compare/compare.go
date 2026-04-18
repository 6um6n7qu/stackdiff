// Package compare provides a high-level pipeline for comparing two config maps
// and producing an annotated, enriched diff report.
package compare

import (
	"github.com/yourorg/stackdiff/internal/diff"
	"github.com/yourorg/stackdiff/internal/enrich"
	"github.com/yourorg/stackdiff/internal/filter"
	"github.com/yourorg/stackdiff/internal/mask"
	"github.com/yourorg/stackdiff/internal/normalize"
	"github.com/yourorg/stackdiff/internal/report"
)

// Options controls pipeline behaviour.
type Options struct {
	FilterStatuses []string
	KeyPrefix      string
	RedactSecrets  bool
	Normalize      bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		RedactSecrets: true,
		Normalize:     true,
	}
}

// Run compares two raw config maps and returns a populated report.
func Run(a, b map[string]string, opts Options) (*report.Report, error) {
	entries := diff.Compare(a, b)

	if opts.Normalize {
		norm := normalize.New(normalize.DefaultConfig())
		entries = norm.Apply(entries)
	}

	if opts.RedactSecrets {
		m := mask.New(nil)
		entries = mask.ApplyToEntries(m, entries)
	}

	if len(opts.FilterStatuses) > 0 || opts.KeyPrefix != "" {
		entries = filter.Apply(entries, filter.Options{
			Statuses:  opts.FilterStatuses,
			KeyPrefix: opts.KeyPrefix,
		})
	}

	en := enrich.New(enrich.DefaultRules())
	entries = en.Apply(entries)

	r := report.New(entries, map[string]string{
		"source": "compare.Run",
	})
	return r, nil
}
