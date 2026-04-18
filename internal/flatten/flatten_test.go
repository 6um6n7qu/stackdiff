package flatten_test

import (
	"testing"

	"github.com/stackdiff/stackdiff/internal/diff"
	"github.com/stackdiff/stackdiff/internal/flatten"
)

func TestFlatten_FlatMap(t *testing.T) {
	m := map[string]any{
		"host": "localhost",
		"port": "8080",
	}
	entries := flatten.Flatten("", m)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	keys := map[string]string{}
	for _, e := range entries {
		keys[e.Key] = e.NewVal
	}
	if keys["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %s", keys["host"])
	}
}

func TestFlatten_NestedMap(t *testing.T) {
	m := map[string]any{
		"db": map[string]any{
			"host": "db.local",
			"port": "5432",
		},
	}
	entries := flatten.Flatten("", m)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "db.host" {
		t.Errorf("expected db.host, got %s", entries[0].Key)
	}
	if entries[1].Key != "db.port" {
		t.Errorf("expected db.port, got %s", entries[1].Key)
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	m := map[string]any{"key": "val"}
	entries := flatten.Flatten("svc", m)
	if entries[0].Key != "svc.key" {
		t.Errorf("expected svc.key, got %s", entries[0].Key)
	}
}

func TestFlatten_StatusIsEqual(t *testing.T) {
	m := map[string]any{"x": "1"}
	entries := flatten.Flatten("", m)
	if entries[0].Status != diff.StatusEqual {
		t.Errorf("expected StatusEqual")
	}
}

func TestFlatten_EmptyMap(t *testing.T) {
	entries := flatten.Flatten("", map[string]any{})
	if len(entries) != 0 {
		t.Errorf("expected 0 entries")
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	m := map[string]any{"a": "1", "b": "2"}
	entries := flatten.Flatten("", m)
	out := flatten.ToMap(entries)
	if out["a"] != "1" || out["b"] != "2" {
		t.Errorf("round-trip mismatch: %v", out)
	}
}
