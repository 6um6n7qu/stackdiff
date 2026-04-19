package diff

import (
	"testing"
	"time"
)

func ts(s string) string { return s }

func entryWithTS(key, timestamp string) Entry {
	e := Entry{Key: key, OldValue: "a", NewValue: "b", Status: StatusChanged}
	if timestamp != "" {
		e.Meta = map[string]string{"timestamp": timestamp}
	}
	return e
}

func TestApplyWindow_NoConfig_ReturnsAll(t *testing.T) {
	entries := []Entry{entryWithTS("a", "2024-01-01T00:00:00Z"), entryWithTS("b", "")}
	out := ApplyWindow(entries, DefaultWindowConfig())
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestApplyWindow_FromFilter(t *testing.T) {
	from, _ := time.Parse(time.RFC3339, "2024-06-01T00:00:00Z")
	cfg := WindowConfig{From: from}
	entries := []Entry{
		entryWithTS("old", "2024-01-01T00:00:00Z"),
		entryWithTS("new", "2024-07-01T00:00:00Z"),
	}
	out := ApplyWindow(entries, cfg)
	if len(out) != 1 || out[0].Key != "new" {
		t.Fatalf("expected only 'new', got %+v", out)
	}
}

func TestApplyWindow_ToFilter(t *testing.T) {
	to, _ := time.Parse(time.RFC3339, "2024-06-01T00:00:00Z")
	cfg := WindowConfig{To: to}
	entries := []Entry{
		entryWithTS("old", "2024-01-01T00:00:00Z"),
		entryWithTS("new", "2024-07-01T00:00:00Z"),
	}
	out := ApplyWindow(entries, cfg)
	if len(out) != 1 || out[0].Key != "old" {
		t.Fatalf("expected only 'old', got %+v", out)
	}
}

func TestApplyWindow_NoTimestamp_AlwaysKept(t *testing.T) {
	from, _ := time.Parse(time.RFC3339, "2024-06-01T00:00:00Z")
	cfg := WindowConfig{From: from}
	e := Entry{Key: "no-ts", Status: StatusAdded}
	out := ApplyWindow([]Entry{e}, cfg)
	if len(out) != 1 {
		t.Fatalf("expected entry without timestamp to be kept")
	}
}

func TestApplyWindow_InvalidTimestamp_Kept(t *testing.T) {
	from, _ := time.Parse(time.RFC3339, "2024-06-01T00:00:00Z")
	cfg := WindowConfig{From: from}
	e := entryWithTS("bad", "not-a-date")
	out := ApplyWindow([]Entry{e}, cfg)
	if len(out) != 1 {
		t.Fatalf("expected entry with bad timestamp to be kept")
	}
}

func TestApplyWindow_BothBounds(t *testing.T) {
	from, _ := time.Parse(time.RFC3339, "2024-03-01T00:00:00Z")
	to, _ := time.Parse(time.RFC3339, "2024-09-01T00:00:00Z")
	cfg := WindowConfig{From: from, To: to}
	entries := []Entry{
		entryWithTS("jan", "2024-01-01T00:00:00Z"),
		entryWithTS("may", "2024-05-01T00:00:00Z"),
		entryWithTS("dec", "2024-12-01T00:00:00Z"),
	}
	out := ApplyWindow(entries, cfg)
	if len(out) != 1 || out[0].Key != "may" {
		t.Fatalf("expected only 'may', got %+v", out)
	}
}
