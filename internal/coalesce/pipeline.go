package coalesce

import (
	"fmt"

	"github.com/yourusername/stackdiff/internal/diff"
)

// Run applies coalescing to the given entries using the default configuration.
// It returns the coalesced entries and any error encountered.
func Run(entries []diff.Entry) ([]diff.Entry, error) {
	cfg := DefaultConfig()
	return Apply(entries, cfg)
}

// RunWithConfig applies coalescing to the given entries using the provided configuration.
func RunWithConfig(entries []diff.Entry, cfg Config) ([]diff.Entry, error) {
	return Apply(entries, cfg)
}

// MustRun applies coalescing using the default configuration and panics on error.
func MustRun(entries []diff.Entry) []diff.Entry {
	out, err := Run(entries)
	if err != nil {
		panic(fmt.Sprintf("coalesce.MustRun: %v", err))
	}
	return out
}

// MustRunWithConfig applies coalescing using the provided configuration and panics on error.
func MustRunWithConfig(entries []diff.Entry, cfg Config) []diff.Entry {
	out, err := RunWithConfig(entries, cfg)
	if err != nil {
		panic(fmt.Sprintf("coalesce.MustRunWithConfig: %v", err))
	}
	return out
}

// Summary holds the result of a coalesce pipeline run.
type Summary struct {
	// Input is the number of entries before coalescing.
	Input int
	// Output is the number of entries after coalescing.
	Output int
	// Coalesced is the number of entries that were merged or dropped.
	Coalesced int
}

// String returns a human-readable summary of the coalesce operation.
func (s Summary) String() string {
	return fmt.Sprintf("coalesce: %d input → %d output (%d coalesced)", s.Input, s.Output, s.Coalesced)
}

// RunWithSummary applies coalescing using the default configuration and returns
// both the resulting entries and a Summary describing what changed.
func RunWithSummary(entries []diff.Entry) ([]diff.Entry, Summary, error) {
	return RunWithConfigAndSummary(entries, DefaultConfig())
}

// RunWithConfigAndSummary applies coalescing using the provided configuration and
// returns both the resulting entries and a Summary describing what changed.
func RunWithConfigAndSummary(entries []diff.Entry, cfg Config) ([]diff.Entry, Summary, error) {
	out, err := Apply(entries, cfg)
	if err != nil {
		return nil, Summary{}, err
	}
	s := Summary{
		Input:     len(entries),
		Output:    len(out),
		Coalesced: Count(entries, cfg),
	}
	return out, s, nil
}
