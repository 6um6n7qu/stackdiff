package patch

import (
	"fmt"
	"io"
	"strings"
)

// PipelineResult holds the generated ops and a human-readable summary.
type PipelineResult struct {
	Ops     []Op
	HasDiff bool
}

// Summary returns a short description of the patch.
func (r PipelineResult) Summary() string {
	var sets, deletes int
	for _, op := range r.Ops {
		if op.Action == "set" {
			sets++
		} else {
			deletes++
		}
	}
	return fmt.Sprintf("%d set, %d delete", sets, deletes)
}

// Run generates a patch from src to dst and writes a human-readable
// representation to w.
func Run(src, dst map[string]string, w io.Writer) (PipelineResult, error) {
	ops := Generate(src, dst)
	result := PipelineResult{Ops: ops, HasDiff: len(ops) > 0}

	if !result.HasDiff {
		_, err := fmt.Fprintln(w, "no changes")
		return result, err
	}

	var sb strings.Builder
	for _, op := range ops {
		sb.WriteString(op.String())
		sb.WriteByte('\n')
	}
	_, err := fmt.Fprint(w, sb.String())
	return result, err
}
