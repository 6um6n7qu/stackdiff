package env

import (
	"testing"
)

func TestResolve_NoPlaceholders(t *testing.T) {
	entries := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	out, results, err := Resolve(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" || out["PORT"] != "8080" {
		t.Errorf("unexpected output: %v", out)
	}
	for _, r := range results {
		if r.Resolved {
			t.Errorf("key %q should not be marked resolved", r.Key)
		}
	}
}

func TestResolve_WithEnvVar(t *testing.T) {
	t.Setenv("APP_SECRET", "supersecret")

	entries := map[string]string{
		"SECRET": "${APP_SECRET}",
	}
	out, results, err := Resolve(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET"] != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", out["SECRET"])
	}
	if len(results) != 1 || !results[0].Resolved {
		t.Errorf("expected result to be marked resolved")
	}
}

func TestResolve_MissingEnvVar(t *testing.T) {
	// Ensure the variable is absent
	t.Setenv("MISSING_VAR", "")
	// Unset it properly by using a key that definitely won't exist
	entries := map[string]string{
		"KEY": "${__STACKDIFF_UNDEFINED_XYZ__}",
	}
	_, _, err := Resolve(entries)
	if err == nil {
		t.Fatal("expected error for missing env var, got nil")
	}
}

func TestResolve_MixedEntries(t *testing.T) {
	t.Setenv("DB_HOST", "db.internal")

	entries := map[string]string{
		"PLAIN":   "value",
		"DB_ADDR": "${DB_HOST}:5432",
	}
	out, results, err := Resolve(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_ADDR"] != "db.internal:5432" {
		t.Errorf("expected 'db.internal:5432', got %q", out["DB_ADDR"])
	}

	resolvedCount := 0
	for _, r := range results {
		if r.Resolved {
			resolvedCount++
		}
	}
	if resolvedCount != 1 {
		t.Errorf("expected 1 resolved entry, got %d", resolvedCount)
	}
}

func TestResolve_EmptyMap(t *testing.T) {
	out, results, err := Resolve(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 || len(results) != 0 {
		t.Errorf("expected empty output, got %v / %v", out, results)
	}
}
