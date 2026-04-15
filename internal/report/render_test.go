package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/stackdiff/internal/diff"
)

func TestRender_Text(t *testing.T) {
	r := New("staging", "prod", sampleDiffs())
	var buf bytes.Buffer
	if err := Render(&buf, r, FormatText); err != nil {
		t.Fatalf("Render text error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[+] APP_ENV") {
		t.Errorf("expected added entry in text output, got:\n%s", out)
	}
	if !strings.Contains(out, "[-] DB_HOST") {
		t.Errorf("expected removed entry in text output, got:\n%s", out)
	}
	if !strings.Contains(out, "[~] LOG_LEVEL") {
		t.Errorf("expected changed entry in text output, got:\n%s", out)
	}
}

func TestRender_JSON(t *testing.T) {
	r := New("staging", "prod", sampleDiffs())
	var buf bytes.Buffer
	if err := Render(&buf, r, FormatJSON); err != nil {
		t.Fatalf("Render JSON error: %v", err)
	}
	var decoded Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("JSON unmarshal error: %v", err)
	}
	if decoded.Summary.Total != 3 {
		t.Errorf("expected total 3 in JSON, got %d", decoded.Summary.Total)
	}
}

func TestRender_Markdown(t *testing.T) {
	r := New("staging", "prod", sampleDiffs())
	var buf bytes.Buffer
	if err := Render(&buf, r, FormatMarkdown); err != nil {
		t.Fatalf("Render markdown error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "# StackDiff Report") {
		t.Errorf("expected markdown header, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected APP_ENV in markdown, got:\n%s", out)
	}
}

func TestRender_DefaultIsText(t *testing.T) {
	r := New("a", "b", []diff.DiffEntry{
		{Key: "X", Kind: diff.Added, NewValue: "1"},
	})
	var buf bytes.Buffer
	if err := Render(&buf, r, Format("unknown")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[+] X") {
		t.Errorf("expected text fallback output")
	}
}
