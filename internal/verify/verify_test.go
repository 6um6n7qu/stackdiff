package verify_test

import (
	"testing"

	"github.com/stackdiff/stackdiff/internal/diff"
	"github.com/stackdiff/stackdiff/internal/verify"
)

func makeEntry(key, old, nw string, status diff.EntryStatus) diff.Entry {
	return diff.Entry{Key: key, OldValue: old, NewValue: nw, Status: status}
}

func TestEvaluate_NoDrift(t *testing.T) {
	entries := []diff.Entry{
		makeEntry("HOST", "a", "a", diff.StatusEqual),
	}
	r := verify.Evaluate(entries, verify.DefaultConfig())
	if !r.Passed() {
		t.Fatalf("expected pass, got %s", r.Status)
	}
}

func TestEvaluate_FailOnChanged(t *testing.T) {
	entries := []diff.Entry{
		makeEntry("PORT", "8080", "9090", diff.StatusChanged),
	}
	r := verify.Evaluate(entries, verify.DefaultConfig())
	if r.Status != verify.StatusFail {
		t.Fatalf("expected fail, got %s", r.Status)
	}
	if len(r.Messages) == 0 {
		t.Fatal("expected messages")
	}
}

func TestEvaluate_IgnoreAdded_Warning(t *testing.T) {
	entries := []diff.Entry{
		makeEntry("NEW_KEY", "", "val", diff.StatusAdded),
	}
	cfg := verify.DefaultConfig()
	cfg.IgnoreAdded = true
	r := verify.Evaluate(entries, cfg)
	if r.Status != verify.StatusWarning {
		t.Fatalf("expected warning, got %s", r.Status)
	}
}

func TestEvaluate_IgnoreRemoved_Warning(t *testing.T) {
	entries := []diff.Entry{
		makeEntry("OLD_KEY", "val", "", diff.StatusRemoved),
	}
	cfg := verify.DefaultConfig()
	cfg.IgnoreRemoved = true
	r := verify.Evaluate(entries, cfg)
	if r.Status != verify.StatusWarning {
		t.Fatalf("expected warning, got %s", r.Status)
	}
}

func TestEvaluate_FailOverridesWarning(t *testing.T) {
	entries := []diff.Entry{
		makeEntry("NEW_KEY", "", "val", diff.StatusAdded),
		makeEntry("PORT", "80", "443", diff.StatusChanged),
	}
	cfg := verify.DefaultConfig()
	cfg.IgnoreAdded = true
	r := verify.Evaluate(entries, cfg)
	if r.Status != verify.StatusFail {
		t.Fatalf("expected fail, got %s", r.Status)
	}
}

func TestResult_Passed(t *testing.T) {
	r := verify.Result{Status: verify.StatusPass}
	if !r.Passed() {
		t.Fatal("expected Passed() true")
	}
	r2 := verify.Result{Status: verify.StatusFail}
	if r2.Passed() {
		t.Fatal("expected Passed() false")
	}
}
