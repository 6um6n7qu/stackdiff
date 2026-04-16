package audit_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stackdiff/internal/audit"
	"github.com/stackdiff/internal/diff"
)

func sampleEvent() audit.Event {
	return audit.Event{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     audit.LevelWarn,
		Message:   "drift detected: 1 change(s)",
		Entries:   []diff.Entry{{Key: "PORT", Left: "8080", Right: "9090", Status: diff.StatusChanged}},
	}
}

func TestParseFormat_Valid(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"text", "text"}, {"json", "json"}, {"", "text"},
	} {
		f, err := audit.ParseFormat(tc.in)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.in, err)
		}
		if string(f) != tc.want {
			t.Errorf("got %q, want %q", f, tc.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := audit.ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestFormatEvent_Text(t *testing.T) {
	s, err := audit.FormatEvent(sampleEvent(), audit.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(s, "WARN") {
		t.Errorf("expected WARN in text output: %s", s)
	}
}

func TestFormatEvent_JSON(t *testing.T) {
	s, err := audit.FormatEvent(sampleEvent(), audit.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(s, "\"Level\"") {
		t.Errorf("expected JSON keys in output: %s", s)
	}
}
