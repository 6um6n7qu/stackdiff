package export

import (
	"fmt"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/filter"
	"github.com/user/stackdiff/internal/report"
	"github.com/user/stackdiff/internal/report/render"
)

// PipelineInput holds everything needed to run an export pipeline.
type PipelineInput struct {
	Entries []diff.Entry
	Labels  map[string]string
	Filter  filter.Options
	Export  Options
}

// RunPipeline filters entries, builds a report, renders it, and writes output.
func RunPipeline(input PipelineInput) error {
	filtered := filter.Apply(input.Entries, input.Filter)

	rep := report.New(filtered, input.Labels)

	exp := New(input.Export)

	w, err := exp.Writer()
	if err != nil {
		return fmt.Errorf("pipeline: %w", err)
	}
	defer w.Close()

	rendOpts := render.Options{
		Format: string(exp.Format()),
		Writer: w,
	}
	if err := render.Render(rep, rendOpts); err != nil {
		return fmt.Errorf("pipeline: render: %w", err)
	}
	return nil
}
