package stale

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"stackdiff/internal/diff"
)

// PipelineInput bundles a snapshot with its capture timestamp.
type PipelineInput struct {
	Entries   []diff.Entry
	CapturedAt time.Time
}

// Run detects stale entries from a pipeline of inputs and writes a summary
// to w. Returns the list of stale results.
func Run(inputs []PipelineInput, cfg Config, w io.Writer) ([]Result, error) {
	if w == nil {
		return nil, fmt.Errorf("stale: writer must not be nil")
	}

	snaps := make([][]diff.Entry, len(inputs))
	timestamps := make([]time.Time, len(inputs))
	for i, inp := range inputs {
		snaps[i] = inp.Entries
		timestamps[i] = inp.CapturedAt
	}

	results := Detect(snaps, timestamps, cfg)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tSTATUS\tREPEATS\tSINCE")
	for _, r := range results {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\n",
			r.Entry.Key,
			r.Entry.Status,
			r.Repeats,
			r.Since.Format(time.RFC3339),
		)
	}
	_ = tw.Flush()
	return results, nil
}
