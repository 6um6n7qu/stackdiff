// Package pipeline wires together the canonical stackdiff comparison pipeline:
// load → compare → chain steps → report.
package pipeline

import (
	"fmt"
	"io"

	"github.com/example/stackdiff/internal/diff"
	"github.com/example/stackdiff/internal/report"
)

// Options controls which optional steps run.
type Options struct {
	Steps  []diff.Step
	Format string // text | json | markdown
	Out    io.Writer
}

// Result is the output of a pipeline run.
type Result struct {
	Entries []diff.Entry
	Report  *report.Report
}

// Run executes the full pipeline given two config maps.
func Run(a, b map[string]string, opts Options) (*Result, error) {
	if opts.Out == nil {
		return nil, fmt.Errorf("pipeline: Out writer must not be nil")
	}

	entries := diff.Compare(a, b)

	if len(opts.Steps) > 0 {
		entries = diff.Chain(entries, opts.Steps)
	}

	rep := report.New(entries, nil)

	fmt := opts.Format
	if fmt == "" {
		fmt = "text"
	}
	if err := report.Render(rep, fmt, opts.Out); err != nil {
		return nil, fmt.Errorf("pipeline: render: %w", err)
	}

	return &Result{Entries: entries, Report: rep}, nil
}
