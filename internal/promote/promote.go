// Package promote provides utilities for promoting a config snapshot
// from one environment to another, generating a diff and patch set.
package promote

import (
	"fmt"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/patch"
	"github.com/user/stackdiff/internal/snapshot"
)

// Result holds the outcome of a promotion operation.
type Result struct {
	// Source is the environment being promoted from.
	Source string
	// Target is the environment being promoted to.
	Target string
	// Entries are the diff entries between source and target.
	Entries []diff.Entry
	// Ops are the patch operations needed to bring target in line with source.
	Ops []patch.Op
}

// HasChanges reports whether the promotion would alter the target.
func (r Result) HasChanges() bool {
	return len(r.Ops) > 0
}

// Summary returns a human-readable one-liner.
func (r Result) Summary() string {
	return fmt.Sprintf("promote %s -> %s: %d op(s)", r.Source, r.Target, len(r.Ops))
}

// Run compares src against dst and produces the patch ops required to
// promote src's configuration into dst.
func Run(src, dst *snapshot.Snapshot) (Result, error) {
	if src == nil {
		return Result{}, fmt.Errorf("promote: source snapshot is nil")
	}
	if dst == nil {
		return Result{}, fmt.Errorf("promote: target snapshot is nil")
	}

	entries := diff.Compare(src.Entries, dst.Entries)
	ops := patch.Generate(dst.Entries, src.Entries)

	return Result{
		Source:  src.Label,
		Target:  dst.Label,
		Entries: entries,
		Ops:     ops,
	}, nil
}
