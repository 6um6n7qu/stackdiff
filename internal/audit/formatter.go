package audit

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Format controls how audit events are serialised.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// ParseFormat converts a string to a Format, returning an error for unknowns.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	}
	return "", fmt.Errorf("audit: unknown format %q", s)
}

// FormatEvent serialises an Event to a string using the given Format.
func FormatEvent(e Event, f Format) (string, error) {
	switch f {
	case FormatJSON:
		b, err := json.Marshal(e)
		if err != nil {
			return "", fmt.Errorf("audit: json marshal: %w", err)
		}
		return string(b), nil
	default:
		return fmt.Sprintf("[%s] %s %s (%d entries)",
			e.Level,
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Message,
			len(e.Entries),
		), nil
	}
}
