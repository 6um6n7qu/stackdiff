package watermark_test

import (
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/watermark"
)

func pipelineEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "HOST", OldValue: "a", NewValue: "b", Status: diff.StatusChanged},
		{Key: "PORT", OldValue: "", NewValue: "8080", Status: diff.StatusAdded},
		{Key: "LOG", OldValue: "info", NewValue: "info", Status: diff.StatusEqual},
	}
}

func TestRun_CountsDriftEntries(t *testing.T) {
	s := watermark.New(tmpPath(t))
	r := watermark.Run(s, "svc", pipelineEntries())
	if r.Count != 2 {
		t.Fatalf("expected count 2, got %d", r.Count)
	}
}

func TestRun_NewPeakOnFirstCall(t *testing.T) {
	s := watermark.New(tmpPath(t))
	r := watermark.Run(s, "svc", pipelineEntries())
	if !r.NewPeak {
		t.Fatal("expected NewPeak=true on first call")
	}
}

func TestRun_NoPeakWhenCountDrops(t *testing.T) {
	s := watermark.New(tmpPath(t))
	s.Record("svc", 10)
	r := watermark.Run(s, "svc", pipelineEntries())
	if r.NewPeak {
		t.Fatal("expected NewPeak=false when count is below stored peak")
	}
}

func TestRun_StringContainsSeries(t *testing.T) {
	s := watermark.New(tmpPath(t))
	r := watermark.Run(s, "my-service", pipelineEntries())
	if got := r.String(); len(got) == 0 {
		t.Fatal("expected non-empty String()")
	}
}

func TestMustNotExceedPeak_OK(t *testing.T) {
	s := watermark.New(tmpPath(t))
	s.Record("svc", 5)
	err := watermark.MustNotExceedPeak(s, "svc", pipelineEntries(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMustNotExceedPeak_Exceeded(t *testing.T) {
	s := watermark.New(tmpPath(t))
	s.Record("svc", 0)
	// tolerance 0 means any new peak beyond 0 should fail
	err := watermark.MustNotExceedPeak(s, "svc", pipelineEntries(), 0)
	if err == nil {
		t.Fatal("expected error when drift exceeds peak beyond tolerance")
	}
}

func TestRun_EmptyEntries_ZeroCount(t *testing.T) {
	s := watermark.New(tmpPath(t))
	r := watermark.Run(s, "svc", []diff.Entry{})
	if r.Count != 0 {
		t.Fatalf("expected count 0 for empty entries, got %d", r.Count)
	}
}
