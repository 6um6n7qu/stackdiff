package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/stackdiff/internal/diff"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Alert represents a drift alert emitted when changes are detected.
type Alert struct {
	Level     Level
	Message   string
	Entries   []diff.Entry
	Timestamp time.Time
}

// Config controls alert behaviour.
type Config struct {
	Level  Level
	Writer io.Writer
}

// DefaultConfig returns a Config writing to stderr at warn level.
func DefaultConfig() Config {
	return Config{
		Level:  LevelWarn,
		Writer: os.Stderr,
	}
}

// Emit writes an alert for the given drift entries if any exist.
func Emit(entries []diff.Entry, cfg Config) *Alert {
	if len(entries) == 0 {
		return nil
	}
	a := &Alert{
		Level:     cfg.Level,
		Message:   fmt.Sprintf("drift detected: %d change(s)", len(entries)),
		Entries:   entries,
		Timestamp: time.Now().UTC(),
	}
	fmt.Fprintf(cfg.Writer, "[%s] %s %s\n", a.Level, a.Timestamp.Format(time.RFC3339), a.Message)
	return a
}
