package transform_test

import (
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/transform"
)

func TestNewPipeline_RunsAllSteps(t *testing.T) {
	entries := []diff.Entry{
		{Key: "HOST", OldVal: " old ", NewVal: " new ", Status: diff.StatusChanged},
	}
	p := transform.NewPipeline(
		transform.New(transform.TrimSpace()),
		transform.New(transform.LowercaseKeys()),
	)
	out := p.Run(entries)
	if out[0].Key != "host" {
		t.Errorf("expected host, got %s", out[0].Key)
	}
	if out[0].OldVal != "old" {
		t.Errorf("expected old, got %q", out[0].OldVal)
	}
}

func TestDefaultPipeline_TrimsAndLowercases(t *testing.T) {
	entries := []diff.Entry{
		{Key: "PORT", OldVal: " 8080 ", NewVal: " 9090 ", Status: diff.StatusChanged},
	}
	p := transform.DefaultPipeline()
	out := p.Run(entries)
	if out[0].Key != "port" {
		t.Errorf("expected port, got %s", out[0].Key)
	}
	if out[0].OldVal != "8080" {
		t.Errorf("expected 8080, got %q", out[0].OldVal)
	}
}

func TestNewPipeline_EmptySteps_ReturnsOriginal(t *testing.T) {
	entries := []diff.Entry{
		{Key: "X", OldVal: "a", NewVal: "b", Status: diff.StatusChanged},
	}
	p := transform.NewPipeline()
	out := p.Run(entries)
	if len(out) != 1 || out[0].Key != "X" {
		t.Error("expected unchanged entries from empty pipeline")
	}
}
