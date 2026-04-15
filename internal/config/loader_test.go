package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLoadFromFile_ValidYAML(t *testing.T) {
	content := "DATABASE_URL: postgres://localhost/dev\nPORT: \"5432\"\nDEBUG: \"true\"\n"
	path := writeTempFile(t, content)

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg["DATABASE_URL"] != "postgres://localhost/dev" {
		t.Errorf("unexpected value for DATABASE_URL: %s", cfg["DATABASE_URL"])
	}
	if cfg["PORT"] != "5432" {
		t.Errorf("unexpected value for PORT: %s", cfg["PORT"])
	}
}

func TestLoadFromFile_MissingFile(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFromFile_InvalidYAML(t *testing.T) {
	content := "KEY: valid\nBAD: [unclosed\n"
	path := writeTempFile(t, content)

	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestLoadFromFile_EmptyFile(t *testing.T) {
	path := writeTempFile(t, "")

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("expected no error for empty file, got: %v", err)
	}
	if len(cfg) != 0 {
		t.Errorf("expected empty config, got %d keys", len(cfg))
	}
}
