package export

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Format represents a supported export format.
type Format string

const (
	FormatJSON     Format = "json"
	FormatText     Format = "text"
	FormatMarkdown Format = "markdown"
)

// ParseFormat parses a string into a Format, returning an error if unsupported.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON, nil
	case "text", "txt":
		return FormatText, nil
	case "markdown", "md":
		return FormatMarkdown, nil
	default:
		return "", fmt.Errorf("unsupported export format: %q (use json, text, or markdown)", s)
	}
}

// Options configures an export operation.
type Options struct {
	Format Format
	Dest   string // file path or "-" for stdout
}

// Exporter writes rendered content to a destination.
type Exporter struct {
	opts Options
}

// New creates a new Exporter with the given options.
func New(opts Options) *Exporter {
	if opts.Format == "" {
		opts.Format = FormatText
	}
	if opts.Dest == "" {
		opts.Dest = "-"
	}
	return &Exporter{opts: opts}
}

// Writer returns an io.WriteCloser for the configured destination.
func (e *Exporter) Writer() (io.WriteCloser, error) {
	if e.opts.Dest == "-" {
		return os.Stdout, nil
	}
	dir := filepath.Dir(e.opts.Dest)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("export: create directories: %w", err)
	}
	f, err := os.Create(e.opts.Dest)
	if err != nil {
		return nil, fmt.Errorf("export: open file %q: %w", e.opts.Dest, err)
	}
	return f, nil
}

// Format returns the configured export format.
func (e *Exporter) Format() Format {
	return e.opts.Format
}

// Dest returns the configured destination path.
func (e *Exporter) Dest() string {
	return e.opts.Dest
}
