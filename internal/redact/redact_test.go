package redact_test

import (
	"testing"

	"github.com/stackdiff/internal/redact"
)

func TestIsSensitive_MatchesPassword(t *testing.T) {
	r := redact.New()
	if !r.IsSensitive("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be sensitive")
	}
}

func TestIsSensitive_MatchesToken(t *testing.T) {
	r := redact.New()
	if !r.IsSensitive("AUTH_TOKEN") {
		t.Error("expected AUTH_TOKEN to be sensitive")
	}
}

func TestIsSensitive_SafeKey(t *testing.T) {
	r := redact.New()
	if r.IsSensitive("APP_PORT") {
		t.Error("expected APP_PORT to not be sensitive")
	}
}

func TestIsSensitive_CaseInsensitive(t *testing.T) {
	r := redact.New()
	if !r.IsSensitive("Api_Key") {
		t.Error("expected Api_Key to be sensitive")
	}
}

func TestApply_MasksSensitive(t *testing.T) {
	r := redact.New()
	input := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"APP_PORT":    "8080",
		"API_TOKEN":   "tok123",
	}
	out := r.Apply(input)
	if out["DB_PASSWORD"] != r.Mask {
		t.Errorf("expected DB_PASSWORD to be masked, got %q", out["DB_PASSWORD"])
	}
	if out["API_TOKEN"] != r.Mask {
		t.Errorf("expected API_TOKEN to be masked, got %q", out["API_TOKEN"])
	}
	if out["APP_PORT"] != "8080" {
		t.Errorf("expected APP_PORT to be unchanged, got %q", out["APP_PORT"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	r := redact.New()
	input := map[string]string{"SECRET_KEY": "abc"}
	r.Apply(input)
	if input["SECRET_KEY"] != "abc" {
		t.Error("Apply must not mutate the input map")
	}
}

func TestApply_EmptyMap(t *testing.T) {
	r := redact.New()
	out := r.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
