package notify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/notify"
)

func driftEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "PORT", OldVal: "8080", NewVal: "9090", Status: diff.StatusChanged},
		{Key: "DEBUG", OldVal: "", NewVal: "true", Status: diff.StatusAdded},
		{Key: "HOST", OldVal: "localhost", NewVal: "localhost", Status: diff.StatusEqual},
	}
}

func TestNotify_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.Config{Channel: notify.ChannelStdout, Writer: &buf})
	entries := []diff.Entry{
		{Key: "HOST", OldVal: "localhost", NewVal: "localhost", Status: diff.StatusEqual},
	}
	if err := n.Notify(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for no drift, got: %q", buf.String())
	}
}

func TestNotify_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.Config{Channel: notify.ChannelStdout, Writer: &buf})
	if err := n.Notify(driftEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 drift") {
		t.Errorf("expected drift count in output, got: %q", out)
	}
}

func TestNotify_MessageContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.Config{Channel: notify.ChannelStdout, Writer: &buf})
	_ = n.Notify(driftEntries())
	out := buf.String()
	for _, key := range []string{"PORT", "DEBUG"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in output, got: %q", key, out)
		}
	}
}

func TestDefaultConfig_IsStdout(t *testing.T) {
	cfg := notify.DefaultConfig()
	if cfg.Channel != notify.ChannelStdout {
		t.Errorf("expected stdout channel, got %q", cfg.Channel)
	}
}

func TestNew_NilWriterFallsBack(t *testing.T) {
	n := notify.New(notify.Config{Channel: notify.ChannelStdout})
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
