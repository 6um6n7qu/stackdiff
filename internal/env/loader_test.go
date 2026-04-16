package env

import (
	"os"
	"testing"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadFile_BasicPairs(t *testing.T) {
	path := writeTempEnvFile(t, "HOST=localhost\nPORT=8080\n")
	m, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["HOST"] != "localhost" || m["PORT"] != "8080" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestLoadFile_IgnoresComments(t *testing.T) {
	path := writeTempEnvFile(t, "# comment\nKEY=value\n")
	m, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m["# comment"]; ok {
		t.Error("comment line should not be parsed as key")
	}
	if m["KEY"] != "value" {
		t.Errorf("expected value, got %q", m["KEY"])
	}
}

func TestLoadFile_QuotedValues(t *testing.T) {
	path := writeTempEnvFile(t, `DSN="user:pass@tcp(localhost)/db"`+"\n")
	m, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["DSN"] != "user:pass@tcp(localhost)/db" {
		t.Errorf("unexpected DSN: %q", m["DSN"])
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := LoadFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFile_EmptyFile(t *testing.T) {
	path := writeTempEnvFile(t, "")
	m, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}
