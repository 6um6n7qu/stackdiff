package expire

import (
	"testing"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

func pipelineEntries(past, future string) []diff.Entry {
	return []diff.Entry{
		makeEntry("live", "1", future),
		makeEntry("dead", "2", past),
		makeEntry("ageless", "3", ""),
	}
}

func TestRun_DropsExpired(t *testing.T) {
	past := time.Now().Add(-time.Hour).Format(time.RFC3339)
	future := time.Now().Add(time.Hour).Format(time.RFC3339)
	res := Run(pipelineEntries(past, future))
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Dropped != 1 {
		t.Fatalf("expected Dropped=1, got %d", res.Dropped)
	}
}

func TestRunWithConfig_CustomKey(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MetaKey = "ttl"
	cfg.Now = func() time.Time { return fixedNow }
	past := fixedNow.Add(-time.Minute).Format(time.RFC3339)
	e := diff.Entry{Key: "x", Meta: map[string]string{"ttl": past}}
	res := RunWithConfig([]diff.Entry{e}, cfg)
	if res.Dropped != 1 {
		t.Fatalf("expected Dropped=1 with custom key, got %d", res.Dropped)
	}
}

func TestMustNoneExpired_OK(t *testing.T) {
	future := time.Now().Add(time.Hour).Format(time.RFC3339)
	entries := []diff.Entry{makeEntry("k", "v", future)}
	if err := MustNoneExpired(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMustNoneExpired_Error(t *testing.T) {
	past := time.Now().Add(-time.Hour).Format(time.RFC3339)
	entries := []diff.Entry{
		makeEntry("a", "1", past),
		makeEntry("b", "2", past),
	}
	err := MustNoneExpired(entries)
	if err == nil {
		t.Fatal("expected error for expired entries")
	}
}

func TestMustNoneExpired_Plural(t *testing.T) {
	past := time.Now().Add(-time.Hour).Format(time.RFC3339)
	entries := []diff.Entry{makeEntry("a", "1", past)}
	err := MustNoneExpired(entries)
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	if len(msg) == 0 {
		t.Fatal("expected non-empty error message")
	}
}
