package diff

import "time"

// Window filters entries to those modified within a given time range.
// Entries without a timestamp are always included.

// WindowConfig holds the time bounds for filtering.
type WindowConfig struct {
	From time.Time
	To   time.Time
}

// DefaultWindowConfig returns an open window (no filtering).
func DefaultWindowConfig() WindowConfig {
	return WindowConfig{
		From: time.Time{},
		To:   time.Time{},
	}
}

// ApplyWindow filters a slice of Entry values to those whose Timestamp
// falls within the window. Entries with a zero Timestamp are always kept.
func ApplyWindow(entries []Entry, cfg WindowConfig) []Entry {
	if cfg.From.IsZero() && cfg.To.IsZero() {
		return entries
	}
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		ts, ok := e.Meta["timestamp"]
		if !ok || ts == "" {
			out = append(out, e)
			continue
		}
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			// unparseable timestamp: include the entry
			out = append(out, e)
			continue
		}
		if !cfg.From.IsZero() && t.Before(cfg.From) {
			continue
		}
		if !cfg.To.IsZero() && t.After(cfg.To) {
			continue
		}
		out = append(out, e)
	}
	return out
}
