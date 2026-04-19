package prune_test

import (
	"testing"
	"time"

	"github.com/stackdiff/internal/diff"
	"github.com/stackdiff/internal/prune"
)

func makeEntry(key, status string, lastSeen *time.Time) diff.Entry {
	e := diff.Entry{Key: key, Status: status, Meta: map[string]any{}}
	if lastSeen != nil {
		e.Meta["last_seen"] = *lastSeen
	}
	return e
}

func TestApply_KeepsRecentEntries(t *testing.T) {
	now := time.Now()
	recent := makeEntry("PORT", diff.StatusChanged, ptr(now.Add(-1*time.Hour)))
	cfg := prune.DefaultConfig()
	res := prune.Apply([]diff.Entry{recent}, cfg, now)
	if len(res.Kept) != 1 || len(res.Pruned) != 0 {
		t.Fatalf("expected 1 kept, 0 pruned; got %d kept, %d pruned", len(res.Kept), len(res.Pruned))
	}
}

func TestApply_PrunesOldEntries(t *testing.T) {
	now := time.Now()
	old := makeEntry("HOST", diff.StatusChanged, ptr(now.Add(-10*24*time.Hour)))
	cfg := prune.DefaultConfig()
	res := prune.Apply([]diff.Entry{old}, cfg, now)
	if len(res.Pruned) != 1 || len(res.Kept) != 0 {
		t.Fatalf("expected 0 kept, 1 pruned; got %d kept, %d pruned", len(res.Kept), len(res.Pruned))
	}
}

func TestApply_NoTimestampAlwaysKept(t *testing.T) {
	now := time.Now()
	e := makeEntry("DB_URL", diff.StatusAdded, nil)
	cfg := prune.DefaultConfig()
	res := prune.Apply([]diff.Entry{e}, cfg, now)
	if len(res.Kept) != 1 {
		t.Fatalf("expected entry without timestamp to be kept")
	}
}

func TestApply_StatusFilter_SkipsNonMatchingStatus(t *testing.T) {
	now := time.Now()
	old := makeEntry("KEY", diff.StatusAdded, ptr(now.Add(-10*24*time.Hour)))
	cfg := prune.DefaultConfig()
	cfg.Statuses = []string{diff.StatusChanged}
	res := prune.Apply([]diff.Entry{old}, cfg, now)
	if len(res.Kept) != 1 {
		t.Fatalf("entry with non-matching status should be kept; got pruned")
	}
}

func TestApply_StatusFilter_PrunesMatchingStatus(t *testing.T) {
	now := time.Now()
	old := makeEntry("KEY", diff.StatusRemoved, ptr(now.Add(-10*24*time.Hour)))
	cfg := prune.DefaultConfig()
	cfg.Statuses = []string{diff.StatusRemoved}
	res := prune.Apply([]diff.Entry{old}, cfg, now)
	if len(res.Pruned) != 1 {
		t.Fatalf("expected matching old entry to be pruned")
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	res := prune.Apply(nil, prune.DefaultConfig(), time.Now())
	if len(res.Kept) != 0 || len(res.Pruned) != 0 {
		t.Fatal("expected empty result for nil input")
	}
}

func ptr(t time.Time) *time.Time { return &t }
