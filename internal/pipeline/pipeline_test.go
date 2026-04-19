package pipeline

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/stackdiff/internal/diff"
)

func TestRun_NoDrift(t *testing.T) {
	a := map[string]string{"HOST": "localhost"}
	b := map[string]string{"HOST": "localhost"}
	var buf bytes.Buffer
	res, err := Run(a, b, Options{Out: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Report.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestRun_DetectsDrift(t *testing.T) {
	a := map[string]string{"HOST": "localhost"}
	b := map[string]string{"HOST": "prod.example.com"}
	var buf bytes.Buffer
	res, err := Run(a, b, Options{Out: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Report.HasDrift() {
		t.Error("expected drift")
	}
}

func TestRun_WithChainStep(t *testing.T) {
	a := map[string]string{"HOST": "localhost", "PORT": "8080"}
	b := map[string]string{"HOST": "prod", "PORT": "8080"}
	filter := diff.Step{
		Name: "only-drift",
		Apply: func(entries []diff.Entry) []diff.Entry {
			var out []diff.Entry
			for _, e := range entries {
				if e.IsDrift() {
					out = append(out, e)
				}
			}
			return out
		},
	}
	var buf bytes.Buffer
	res, err := Run(a, b, Options{Steps: []diff.Step{filter}, Out: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 1 {
		t.Errorf("expected 1 entry after filter, got %d", len(res.Entries))
	}
}

func TestRun_JSONFormat(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	var buf bytes.Buffer
	_, err := Run(a, b, Options{Format: "json", Out: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "{") {
		t.Error("expected JSON output")
	}
}

func TestRun_NilWriterReturnsError(t *testing.T) {
	a := map[string]string{}
	b := map[string]string{}
	_, err := Run(a, b, Options{})
	if err == nil {
		t.Error("expected error for nil writer")
	}
}
