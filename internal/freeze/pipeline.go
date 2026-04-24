package freeze

import (
	"errors"
	"fmt"
	"strings"

	"github.com/user/stackdiff/internal/diff"
)

// Result holds the outcome of a freeze pipeline check.
type Result struct {
	Violations []Violation
}

// OK returns true when no frozen keys have drifted.
func (r Result) OK() bool {
	return len(r.Violations) == 0
}

// Summary returns a multi-line string listing all violations.
func (r Result) Summary() string {
	if r.OK() {
		return "no frozen-key violations"
	}
	lines := make([]string, 0, len(r.Violations))
	for _, v := range r.Violations {
		lines = append(lines, "  "+v.String())
	}
	return fmt.Sprintf("%d frozen-key violation(s):\n%s", len(r.Violations), strings.Join(lines, "\n"))
}

// Run executes the freeze check and returns a Result.
func Run(cfg Config, entries []diff.Entry) Result {
	return Result{Violations: Enforce(cfg, entries)}
}

// MustPass returns an error if any frozen-key violations are found.
func MustPass(cfg Config, entries []diff.Entry) error {
	r := Run(cfg, entries)
	if !r.OK() {
		return errors.New(r.Summary())
	}
	return nil
}
