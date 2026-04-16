package profile

import (
	"os"
	"testing"
)

func TestResolve_FromFlag(t *testing.T) {
	res := Resolve("staging")
	if res.Name != "staging" || res.Source != "flag" {
		t.Errorf("expected flag source staging, got %+v", res)
	}
}

func TestResolve_FromEnv(t *testing.T) {
	os.Setenv("STACKDIFF_PROFILE", "production")
	defer os.Unsetenv("STACKDIFF_PROFILE")

	res := Resolve("")
	if res.Name != "production" || res.Source != "env" {
		t.Errorf("expected env source production, got %+v", res)
	}
}

func TestResolve_Default(t *testing.T) {
	os.Unsetenv("STACKDIFF_PROFILE")
	res := Resolve("")
	if res.Name != "default" || res.Source != "default" {
		t.Errorf("expected default, got %+v", res)
	}
}

func TestResolve_FlagTakesPrecedenceOverEnv(t *testing.T) {
	os.Setenv("STACKDIFF_PROFILE", "production")
	defer os.Unsetenv("STACKDIFF_PROFILE")

	res := Resolve("dev")
	if res.Name != "dev" || res.Source != "flag" {
		t.Errorf("expected flag to win, got %+v", res)
	}
}

func TestResolve_NormalizesCase(t *testing.T) {
	res := Resolve("  Staging  ")
	if res.Name != "staging" {
		t.Errorf("expected normalized name, got %q", res.Name)
	}
}

func TestResolveFromStore_Found(t *testing.T) {
	dir := tmpStore(t)
	s := NewStore(dir)
	p := sampleProfile("qa")
	_ = s.Save(p)

	loaded, res, err := ResolveFromStore(s, "qa")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.Name != "qa" || res.Source != "flag" {
		t.Errorf("unexpected result: %+v / %+v", loaded, res)
	}
}

func TestResolveFromStore_Missing(t *testing.T) {
	dir := tmpStore(t)
	s := NewStore(dir)

	_, _, err := ResolveFromStore(s, "ghost")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}
