package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stackdiff/stackdiff/internal/profile"
)

func tmpStore(t *testing.T) *profile.Store {
	t.Helper()
	dir := t.TempDir()
	return profile.NewStore(dir)
}

func sampleProfile(name string) *profile.Profile {
	return &profile.Profile{
		Name:    name,
		Desc:    "test profile",
		EnvFile: "envs/" + name + ".env",
		Labels:  map[string]string{"env": name},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	s := tmpStore(t)
	p := sampleProfile("staging")
	if err := s.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load("staging")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Name != p.Name || got.Desc != p.Desc || got.EnvFile != p.EnvFile {
		t.Errorf("mismatch: got %+v want %+v", got, p)
	}
}

func TestLoad_MissingProfile(t *testing.T) {
	s := tmpStore(t)
	_, err := s.Load("ghost")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestList_Empty(t *testing.T) {
	s := tmpStore(t)
	names, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestList_MultipleProfiles(t *testing.T) {
	s := tmpStore(t)
	for _, n := range []string{"dev", "staging", "prod"} {
		if err := s.Save(sampleProfile(n)); err != nil {
			t.Fatalf("Save %s: %v", n, err)
		}
	}
	names, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 profiles, got %d: %v", len(names), names)
	}
}

func TestSave_InvalidDir(t *testing.T) {
	f, _ := os.CreateTemp("", "notadir")
	f.Close()
	s := profile.NewStore(filepath.Join(f.Name(), "sub"))
	err := s.Save(sampleProfile("x"))
	if err == nil {
		t.Fatal("expected error saving to invalid dir")
	}
}
