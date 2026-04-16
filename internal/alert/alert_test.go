package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stackdiff/internal/alert"
	"github.com/stackdiff/internal/diff"
)

func driftEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "PORT", OldValue: "8080", NewValue: "9090", Status: diff.StatusChanged},
		{Key: "DEBUG", OldValue: "", NewValue: "true", Status: diff.StatusAdded},
	}
}

func TestEmit_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	cfg := alert.Config{Level: alert.LevelWarn, Writer: &buf}
	a := alert.Emit(nil, cfg)
	if a != nil {
		t.Errorf("expected nil alert for empty entries")
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestEmit_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	cfg := alert.Config{Level: alert.LevelWarn, Writer: &buf}
	a := alert.Emit(driftEntries(), cfg)
	if a == nil {
		t.Fatal("expected non-nil alert")
	}
	if a.Level != alert.LevelWarn {
		t.Errorf("expected level warn, got %s", a.Level)
	}
	if len(a.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(a.Entries))
	}
	if !strings.Contains(buf.String(), "drift detected: 2") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestEmit_MessageContent(t *testing.T) {
	var buf bytes.Buffer
	cfg := alert.Config{Level: alert.LevelError, Writer: &buf}
	alert.Emit(driftEntries(), cfg)
	out := buf.String()
	if !strings.Contains(out, string(alert.LevelError)) {
		t.Errorf("expected level in output, got %q", out)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := alert.DefaultConfig()
	if cfg.Level != alert.LevelWarn {
		t.Errorf("expected default level warn, got %s", cfg.Level)
	}
	if cfg.Writer == nil {
		t.Error("expected non-nil writer")
	}
}
