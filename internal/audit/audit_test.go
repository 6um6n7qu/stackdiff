package audit_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stackdiff/internal/audit"
	"github.com/stackdiff/internal/diff"
)

func driftEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "DB_HOST", Left: "localhost", Right: "prod-db", Status: diff.StatusChanged},
		{Key: "NEW_KEY", Left: "", Right: "value", Status: diff.StatusAdded},
	}
}

func TestLog_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	l.Log(audit.LevelInfo, nil, "test message")
	if !strings.Contains(buf.String(), "INFO") {
		t.Errorf("expected INFO in output, got: %s", buf.String())
	}
}

func TestLogDrift_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	e := l.LogDrift(nil)
	if e.Level != audit.LevelInfo {
		t.Errorf("expected INFO, got %s", e.Level)
	}
}

func TestLogDrift_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	e := l.LogDrift(driftEntries())
	if e.Level != audit.LevelWarn {
		t.Errorf("expected WARN, got %s", e.Level)
	}
	if len(e.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(e.Entries))
	}
}

func TestNew_NilWriterUsesStdout(t *testing.T) {
	// Should not panic
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}
