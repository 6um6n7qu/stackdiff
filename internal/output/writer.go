package output

import (
	"fmt"
	"io"
	"os"
)

// Destination represents where output should be written.
type Destination int

const (
	Stdout Destination = iota
	File
)

// Writer wraps an io.Writer with destination metadata.
type Writer struct {
	w    io.Writer
	dest Destination
	path string
}

// NewStdout returns a Writer that writes to standard output.
func NewStdout() *Writer {
	return &Writer{
		w:    os.Stdout,
		dest: Stdout,
	}
}

// NewFile returns a Writer that writes to the given file path.
// The file is created or truncated on first write.
func NewFile(path string) (*Writer, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("output: cannot open file %q: %w", path, err)
	}
	return &Writer{
		w:    f,
		dest: File,
		path: path,
	}, nil
}

// Write writes p to the underlying writer.
func (w *Writer) Write(p []byte) (int, error) {
	return w.w.Write(p)
}

// WriteString writes s to the underlying writer.
func (w *Writer) WriteString(s string) (int, error) {
	return fmt.Fprint(w.w, s)
}

// Destination returns the destination type of this writer.
func (w *Writer) Destination() Destination {
	return w.dest
}

// Path returns the file path if the destination is File, otherwise empty string.
func (w *Writer) Path() string {
	return w.path
}

// Close closes the underlying writer if it implements io.Closer.
func (w *Writer) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
