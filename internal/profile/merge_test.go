package profile

import (
	"testing"
)

func baseProfile() Profile {
	return Profile{
		Name: "base",
		Env:  map[string]string{"HOST": "localhost", "PORT": "8080", "DEBUG": "false"},
	}
}

func overlayProfile() Profile {
	return Profile{
		Name: "overlay",
		Env:  map[string]string{"PORT": "9090", "LOG": "info"},
	}
}

func TestMergeProfiles_NoConflict(t *testing.T) {
	base := Profile{Name: "a", Env: map[string]string{"A": "1"}}
	overlay := Profile{Name: "b", Env: map[string]string{"B": "2"}}
	res, err := MergeProfiles(base, overlay, StrategyLeft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["A"] != "1" || res.Merged["B"] != "2" {
		t.Errorf("unexpected merged map: %v", res.Merged)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestMergeProfiles_StrategyLeft(t *testing.T) {
	res, err := MergeProfiles(baseProfile(), overlayProfile(), StrategyLeft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", res.Merged["PORT"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "PORT" {
		t.Errorf("expected conflict on PORT, got %v", res.Conflicts)
	}
}

func TestMergeProfiles_StrategyRight(t *testing.T) {
	res, err := MergeProfiles(baseProfile(), overlayProfile(), StrategyRight)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["PORT"] != "9090" {
		t.Errorf("expected PORT=9090, got %s", res.Merged["PORT"])
	}
}

func TestMergeProfiles_StrategyError(t *testing.T) {
	_, err := MergeProfiles(baseProfile(), overlayProfile(), StrategyError)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMergeProfiles_SameValue_NoConflict(t *testing.T) {
	base := Profile{Name: "a", Env: map[string]string{"X": "same"}}
	overlay := Profile{Name: "b", Env: map[string]string{"X": "same"}}
	res, err := MergeProfiles(base, overlay, StrategyError)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts for same value, got %v", res.Conflicts)
	}
}
