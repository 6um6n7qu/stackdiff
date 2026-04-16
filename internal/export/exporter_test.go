package export_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/stackdiff/internal/export"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  export.Format
	}{
		{"json", export.FormatJSON},
		{"JSON", export.FormatJSON},
		{"text", export.FormatText},
		{"txt", export.FormatText},
		{"markdown", export.FormatMarkdown},
		{"md", export.FormatMarkdown},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := export.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := export.ParseFormat("csv")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestNew_Defaults(t *testing.T) {
	e := export.New(export.Options{})
	if e.Format() != export.FormatText {
		t.Errorf("default format: got %q, want %q", e.Format(), export.FormatText)
	}
	if e.Dest() != "-" {
		t.Errorf("default dest: got %q, want %q", e.Dest(), "-")
	}
}

func TestWriter_Stdout(t *testing.T) {
	e := export.New(export.Options{Dest: "-"})
	w, err := e.Writer()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w != os.Stdout {
		t.Error("expected os.Stdout")
	}
}

func TestWriter_File(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "sub", "out.txt")
	e := export.New(export.Options{Dest: dest})
	w, err := e.Writer()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()
	if _, err := os.Stat(dest); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestWriter_InvalidPath(t *testing.T) {
	e := export.New(export.Options{Dest: "/dev/null/impossible/path.txt"})
	_, err := e.Writer()
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
