package promote_test

import (
	"testing"
	"time"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/promote"
	"github.com/user/stackdiff/internal/snapshot"
)

func makeSnap(label string, pairs map[string]string) *snapshot.Snapshot {
	entries := make([]diff.Entry, 0, len(pairs))
	for k, v := range pairs {
		entries = append(entries, diff.Entry{Key: k, NewValue: v, Status: diff.StatusEqual})
	}
	return &snapshot.Snapshot{Label: label, Entries: entries, CreatedAt: time.Now()}
}

func TestRun_NilSource(t *testing.T) {
	dst := makeSnap("prod", map[string]string{"A": "1"})
	_, err := promote.Run(nil, dst)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestRun_NilTarget(t *testing.T) {
	src := makeSnap("staging", map[string]string{"A": "1"})
	_, err := promote.Run(src, nil)
	if err == nil {
		t.Fatal("expected error for nil target")
	}
}

func TestRun_NoChanges(t *testing.T) {
	src := makeSnap("staging", map[string]string{"A": "1", "B": "2"})
	dst := makeSnap("prod", map[string]string{"A": "1", "B": "2"})
	res, err := promote.Run(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HasChanges() {
		t.Errorf("expected no changes, got %d op(s)", len(res.Ops))
	}
}

func TestRun_DetectsChanges(t *testing.T) {
	src := makeSnap("staging", map[string]string{"A": "new", "C": "3"})
	dst := makeSnap("prod", map[string]string{"A": "old", "B": "2"})
	res, err := promote.Run(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.HasChanges() {
		t.Error("expected changes but got none")
	}
}

func TestRun_SummaryFormat(t *testing.T) {
	src := makeSnap("staging", map[string]string{"X": "1"})
	dst := makeSnap("prod", map[string]string{"X": "2"})
	res, _ := promote.Run(src, dst)
	s := res.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestRunWithConfig_DryRun_DoesNotMutateDst(t *testing.T) {
	src := makeSnap("staging", map[string]string{"KEY": "new"})
	dst := makeSnap("prod", map[string]string{"KEY": "old"})
	origLen := len(dst.Entries)

	cfg := promote.NewConfig(promote.WithDryRun(true))
	_, err := promote.RunWithConfig(src, dst, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dst.Entries) != origLen {
		t.Errorf("dry-run must not mutate dst entries")
	}
}

func TestRunWithConfig_IgnoreKeys(t *testing.T) {
	src := makeSnap("staging", map[string]string{"A": "1", "SECRET": "s"})
	dst := makeSnap("prod", map[string]string{"A": "2", "SECRET": "old"})

	cfg := promote.NewConfig(
		promote.WithDryRun(true),
		promote.WithIgnoreKeys("SECRET"),
	)
	res, err := promote.RunWithConfig(src, dst, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, op := range res.Ops {
		if op.Key == "SECRET" {
			t.Error("SECRET should have been excluded from ops")
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := promote.DefaultConfig()
	if cfg.DryRun {
		t.Error("DryRun should default to false")
	}
	if cfg.IgnoreKeys == nil {
		t.Error("IgnoreKeys should be initialised")
	}
}
