package output_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/stackdiff/internal/output"
)

func TestNewStdout_Destination(t *testing.T) {
	w := output.NewStdout()
	if w.Destination() != output.Stdout {
		t.Errorf("expected Stdout destination, got %v", w.Destination())
	}
	if w.Path() != "" {
		t.Errorf("expected empty path for stdout writer, got %q", w.Path())
	}
}

func TestNewStdout_Write(t *testing.T) {
	w := output.NewStdout()
	n, err := w.WriteString("hello")
	if err != nil {
		t.Fatalf("unexpected error writing to stdout writer: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5 bytes written, got %d", n)
	}
}

func TestNewFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	w, err := output.NewFile(path)
	if err != nil {
		t.Fatalf("unexpected error creating file writer: %v", err)
	}
	defer w.Close()

	if w.Destination() != output.File {
		t.Errorf("expected File destination, got %v", w.Destination())
	}
	if w.Path() != path {
		t.Errorf("expected path %q, got %q", path, w.Path())
	}
}

func TestNewFile_WritesContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	w, err := output.NewFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = w.WriteString("stackdiff output")
	if err != nil {
		t.Fatalf("write error: %v", err)
	}
	w.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if string(data) != "stackdiff output" {
		t.Errorf("expected %q, got %q", "stackdiff output", string(data))
	}
}

func TestNewFile_InvalidPath(t *testing.T) {
	_, err := output.NewFile("/nonexistent/dir/out.txt")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestStdout_CloseIsNoop(t *testing.T) {
	w := output.NewStdout()
	if err := w.Close(); err != nil {
		t.Errorf("expected nil error closing stdout writer, got %v", err)
	}
}
