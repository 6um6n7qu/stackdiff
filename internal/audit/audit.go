package audit

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/stackdiff/internal/diff"
)

// Level represents the severity of an audit event.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Entries   []diff.Entry
}

// Logger writes audit events to a destination.
type Logger struct {
	out io.Writer
}

// New creates a Logger writing to w. Pass nil to use os.Stdout.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{out: w}
}

// Log writes an audit event for the given diff entries.
func (l *Logger) Log(level Level, entries []diff.Entry, msg string) Event {
	e := Event{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Message:   msg,
		Entries:   entries,
	}
	fmt.Fprintf(l.out, "[%s] %s %s (%d entries)\n",
		e.Level, e.Timestamp.Format(time.RFC3339), e.Message, len(e.Entries))
	return e
}

// LogDrift is a convenience wrapper that logs a WARN when drift is present.
func (l *Logger) LogDrift(entries []diff.Entry) Event {
	if len(entries) == 0 {
		return l.Log(LevelInfo, entries, "no drift detected")
	}
	return l.Log(LevelWarn, entries, fmt.Sprintf("drift detected: %d change(s)", len(entries)))
}
