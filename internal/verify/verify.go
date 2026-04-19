// Package verify checks whether two config snapshots are equivalent
// within a configurable tolerance, returning a structured result.
package verify

import (
	"fmt"

	"github.com/stackdiff/stackdiff/internal/diff"
)

// Status represents the outcome of a verification.
type Status string

const (
	StatusPass    Status = "pass"
	StatusFail    Status = "fail"
	StatusWarning Status = "warning"
)

// Config controls verification behaviour.
type Config struct {
	// MaxDrift is the maximum number of drift entries allowed before failing.
	MaxDrift int
	// IgnoreAdded treats added keys as warnings rather than failures.
	IgnoreAdded bool
	// IgnoreRemoved treats removed keys as warnings rather than failures.
	IgnoreRemoved bool
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{MaxDrift: 0}
}

// Result holds the outcome of a verification run.
type Result struct {
	Status   Status
	Messages []string
	Entries  []diff.Entry
}

// Passed returns true when status is pass.
func (r Result) Passed() bool { return r.Status == StatusPass }

// Evaluate verifies entries against cfg and returns a Result.
func Evaluate(entries []diff.Entry, cfg Config) Result {
	var drifted []diff.Entry
	var messages []string

	for _, e := range entries {
		if e.IsDrift() {
			drifted = append(drifted, e)
		}
	}

	if len(drifted) == 0 {
		return Result{Status: StatusPass, Entries: drifted}
	}

	status := StatusFail
	for _, e := range drifted {
		switch {
		case e.Status == diff.StatusAdded && cfg.IgnoreAdded:
			messages = append(messages, fmt.Sprintf("warning: added key %q", e.Key))
			if status == StatusFail {
				status = StatusWarning
			}
		case e.Status == diff.StatusRemoved && cfg.IgnoreRemoved:
			messages = append(messages, fmt.Sprintf("warning: removed key %q", e.Key))
			if status == StatusFail {
				status = StatusWarning
			}
		default:
			messages = append(messages, fmt.Sprintf("fail: %s key %q", e.Status, e.Key))
			status = StatusFail
		}
	}

	if cfg.MaxDrift > 0 && len(drifted) <= cfg.MaxDrift && status != StatusFail {
		status = StatusWarning
	}

	return Result{Status: status, Messages: messages, Entries: drifted}
}
