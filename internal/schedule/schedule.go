// Package schedule provides periodic diff execution support.
package schedule

import (
	"context"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// Job holds configuration for a scheduled diff run.
type Job struct {
	Interval time.Duration
	Loader   func() (map[string]string, map[string]string, error)
	OnDrift  func([]diff.Entry)
	OnError  func(error)
}

// Run starts the scheduled job and blocks until ctx is cancelled.
func Run(ctx context.Context, job Job) {
	if job.Interval <= 0 {
		job.Interval = time.Minute
	}
	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			execute(job)
		}
	}
}

{
	a, b, err := job. != nil {
		job.OnError(err)
		}
		return
	}
	entries := diff.Compare(a, b)
	var drifted []diff.Entry
	for _, e := range entries {
		if e.IsDrift() {
			drifted = append(drifted, e)
		}
	}
	if len(drifted) > 0 && job.OnDrift != nil {
		job.OnDrift(drifted)
	}
}
