package promote

import (
	"fmt"

	"github.com/user/stackdiff/internal/patch"
	"github.com/user/stackdiff/internal/snapshot"
)

// RunWithConfig executes a promotion using the supplied Config.
// When DryRun is true the patch ops are returned but patch.Apply is not called.
// Keys listed in Config.IgnoreKeys are stripped from the ops before application.
func RunWithConfig(src, dst *snapshot.Snapshot, cfg Config) (Result, error) {
	res, err := Run(src, dst)
	if err != nil {
		return Result{}, err
	}

	// Filter out ignored keys.
	filtered := make([]patch.Op, 0, len(res.Ops))
	for _, op := range res.Ops {
		if _, skip := cfg.IgnoreKeys[op.Key]; skip {
			continue
		}
		filtered = append(filtered, op)
	}
	res.Ops = filtered

	if cfg.DryRun {
		return res, nil
	}

	applied, err := patch.Apply(dst.Entries, res.Ops)
	if err != nil {
		return res, fmt.Errorf("promote: apply failed: %w", err)
	}
	dst.Entries = applied
	return res, nil
}
