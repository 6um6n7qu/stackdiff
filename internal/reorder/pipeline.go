package reorder

import (
	"fmt"

	"github.com/stackdiff/stackdiff/internal/diff"
)

// Run applies the default reorder strategy to entries and returns the result.
func Run(entries []diff.Entry) []diff.Entry {
	return Apply(entries, DefaultConfig())
}

// RunWithStrategy applies the named strategy to entries.
// Returns an error if the strategy is unrecognised.
func RunWithStrategy(entries []diff.Entry, s Strategy) ([]diff.Entry, error) {
	switch s {
	case ByKey, ByKeyDesc, ByStatus:
		return Apply(entries, Config{Strategy: s, StableEqual: true}), nil
	default:
		return nil, fmt.Errorf("reorder: unknown strategy %q", s)
	}
}

// MustRun applies the named strategy and panics on an unrecognised strategy.
// Useful in init paths where the strategy is a compile-time constant.
func MustRun(entries []diff.Entry, s Strategy) []diff.Entry {
	out, err := RunWithStrategy(entries, s)
	if err != nil {
		panic(err)
	}
	return out
}
