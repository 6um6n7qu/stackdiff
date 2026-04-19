package replay

import (
	"fmt"
	"io"

	"github.com/user/stackdiff/internal/history"
	"github.com/user/stackdiff/internal/report"
	"github.com/user/stackdiff/internal/report/render"
)

// PipelineConfig holds output settings for a replay pipeline run.
type PipelineConfig struct {
	Format string
	Writer io.Writer
	Labels map[string]string
}

// RunPipeline replays records from store, renders each as a report, and
// writes the output to cfg.Writer using cfg.Format.
func RunPipeline(store *history.Store, cfg PipelineConfig, opts ...Option) error {
	if cfg.Writer == nil {
		return fmt.Errorf("replay pipeline: writer must not be nil")
	}
	fmt := cfg.Format
	if fmt == "" {
		fmt = "text"
	}

	return Run(store, func(r history.Record) error {
		rpt := report.New(r.Entries, report.WithLabels(cfg.Labels))
		out, err := render.Render(rpt, fmt)
		if err != nil {
			return fmt.Errorf("replay pipeline: render: %w", err)
		}
		_, werr := fmt.Fprintln(cfg.Writer, out)
		return werr
	}, opts...)
}
