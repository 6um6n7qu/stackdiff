package config

import (
	"errors"
	"testing"
)

func TestValidate_ValidConfig(t *testing.T) {
	cfg := map[string]string{
		"DATABASE_URL": "postgres://localhost/dev",
		"PORT":         "5432",
	}
	if err := Validate(cfg); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_NilConfig(t *testing.T) {
	err := Validate(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestValidate_EmptyConfig(t *testing.T) {
	err := Validate(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty config")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) == 0 {
		t.Error("expected at least one issue")
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	cfg := map[string]string{
		"PORT":    "",
		"API_KEY": "secret",
	}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for empty value")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(ve.Issues))
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	cfg := map[string]string{
		"KEY_A": "",
		"KEY_B": "",
	}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(ve.Issues))
	}
}
