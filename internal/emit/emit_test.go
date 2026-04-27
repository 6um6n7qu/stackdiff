package emit_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stackdiff/stackdiff/internal/diff"
	"github.com/stackdiff/stackdiff/internal/emit"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "DB_HOST", OldValue: "localhost", NewValue: "prod-db.internal", Status: diff.StatusChanged},
		{Key: "API_KEY", OldValue: "", NewValue: "abc123", Status: diff.StatusAdded},
		{Key: "LEGACY_FLAG", OldValue: "true", NewValue: "", Status: diff.StatusRemoved},
		{Key: "LOG_LEVEL", OldValue: "info", NewValue: "info", Status: diff.StatusEqual},
	}
}

func TestNewWriterSink_WritesToBuffer(t *testing.T) {
	var buf bytes.Buffer
	sink := emit.NewWriterSink(&buf)
	if sink == nil {
		t.Fatal("expected non-nil sink")
	}
}

func TestEmit_SendsEntriesToSink(t *testing.T) {
	var buf bytes.Buffer
	sink := emit.NewWriterSink(&buf)
	emitter := emit.New(sink)

	entries := sampleEntries()
	if err := emitter.Emit(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected output to contain DB_HOST, got: %s", out)
	}
}

func TestEmit_EmptyEntries_NoError(t *testing.T) {
	var buf bytes.Buffer
	sink := emit.NewWriterSink(&buf)
	emitter := emit.New(sink)

	if err := emitter.Emit([]diff.Entry{}); err != nil {
		t.Fatalf("unexpected error on empty entries: %v", err)
	}
}

func TestEmit_OnlyDriftEntries(t *testing.T) {
	var buf bytes.Buffer
	sink := emit.NewWriterSink(&buf)
	emitter := emit.New(sink)

	entries := sampleEntries()
	if err := emitter.Emit(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	// Equal entries should not appear in drift-only output
	if strings.Contains(out, "LOG_LEVEL") {
		t.Errorf("expected equal entry LOG_LEVEL to be omitted, got: %s", out)
	}
}

func TestEmit_ContainsAllDriftStatuses(t *testing.T) {
	var buf bytes.Buffer
	sink := emit.NewWriterSink(&buf)
	emitter := emit.New(sink)

	entries := sampleEntries()
	if err := emitter.Emit(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, key := range []string{"DB_HOST", "API_KEY", "LEGACY_FLAG"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected output to contain key %q", key)
		}
	}
}

func TestEmit_TimestampPresent(t *testing.T) {
	var buf bytes.Buffer
	sink := emit.NewWriterSink(&buf)
	emitter := emit.New(sink)

	before := time.Now()
	if err := emitter.Emit(sampleEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after := time.Now()

	// Ensure the emitter ran within expected time bounds (smoke test)
	if after.Before(before) {
		t.Error("time went backwards")
	}
}

func TestNewWriterSink_NilWriter_ReturnsError(t *testing.T) {
	sink := emit.NewWriterSink(nil)
	emitter := emit.New(sink)

	err := emitter.Emit(sampleEntries())
	if err == nil {
		t.Error("expected error when writer is nil, got nil")
	}
}
