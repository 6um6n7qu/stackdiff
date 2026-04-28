package watermark_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/stackdiff/internal/watermark"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "watermark.json")
}

func TestRecord_NewPeak(t *testing.T) {
	s := watermark.New(tmpPath(t))
	if !s.Record("svc", 5) {
		t.Fatal("expected new peak on first record")
	}
}

func TestRecord_NoPeakWhenLower(t *testing.T) {
	s := watermark.New(tmpPath(t))
	s.Record("svc", 10)
	if s.Record("svc", 3) {
		t.Fatal("expected no new peak when count is lower")
	}
}

func TestRecord_PeakUpdatedWhenHigher(t *testing.T) {
	s := watermark.New(tmpPath(t))
	s.Record("svc", 5)
	if !s.Record("svc", 8) {
		t.Fatal("expected new peak when count is higher")
	}
	m, _ := s.Get("svc")
	if m.Peak != 8 {
		t.Fatalf("expected peak 8, got %d", m.Peak)
	}
}

func TestGet_MissingSeriesReturnsFalse(t *testing.T) {
	s := watermark.New(tmpPath(t))
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected false for unknown series")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	p := tmpPath(t)
	s := watermark.New(p)
	s.Record("alpha", 7)
	if err := s.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s2 := watermark.New(p)
	if err := s2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	m, ok := s2.Get("alpha")
	if !ok || m.Peak != 7 {
		t.Fatalf("expected peak 7 after round-trip, got %+v", m)
	}
}

func TestLoad_MissingFile_NoError(t *testing.T) {
	s := watermark.New(filepath.Join(t.TempDir(), "nope.json"))
	if err := s.Load(); err != nil {
		t.Fatalf("unexpected error loading missing file: %v", err)
	}
}

func TestSave_InvalidPath(t *testing.T) {
	s := watermark.New("/no/such/dir/watermark.json")
	if err := s.Save(); err == nil {
		t.Fatal("expected error saving to invalid path")
	}
}

func TestReset_ClearsMark(t *testing.T) {
	s := watermark.New(tmpPath(t))
	s.Record("svc", 5)
	s.Reset("svc")
	_, ok := s.Get("svc")
	if ok {
		t.Fatal("expected mark to be cleared after Reset")
	}
}

func TestLoad_IgnoresEnvFile(t *testing.T) {
	p := tmpPath(t)
	if err := os.WriteFile(p, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := watermark.New(p)
	if err := s.Load(); err == nil {
		t.Fatal("expected error on invalid JSON")
	}
}
