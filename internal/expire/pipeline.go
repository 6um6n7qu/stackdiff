package expire

import (
	"fmt"

	"github.com/user/stackdiff/internal/diff"
)

// PipelineResult holds the output of a pipeline run.
type PipelineResult struct {
	Entries []diff.Entry
	Dropped int
}

// Run applies expiration with the default config and returns a PipelineResult.
func Run(entries []diff.Entry) PipelineResult {
	return RunWithConfig(entries, DefaultConfig())
}

// RunWithConfig applies expiration with a custom config.
func RunWithConfig(entries []diff.Entry, cfg Config) PipelineResult {
	dropped := CountExpired(entries, cfg)
	filtered := Apply(entries, cfg)
	return PipelineResult{
		Entries: filtered,
		Dropped: dropped,
	}
}

// MustNoneExpired returns an error if any entries are expired under the
// default config. Useful as a gate in stricter pipelines.
func MustNoneExpired(entries []diff.Entry) error {
	return MustNoneExpiredWithConfig(entries, DefaultConfig())
}

// MustNoneExpiredWithConfig is like MustNoneExpired but accepts a custom Config.
func MustNoneExpiredWithConfig(entries []diff.Entry, cfg Config) error {
	n := CountExpired(entries, cfg)
	if n > 0 {
		return fmt.Errorf("expire: %d expired entr%s detected", n, plural(n))
	}
	return nil
}

func plural(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}
