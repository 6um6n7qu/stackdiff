package env

import (
	"testing"
)

func TestMergeMaps_NoConflict(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"C": "3"}
	res, err := MergeMaps(left, right, StrategyLeft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Values) != 3 {
		t.Errorf("expected 3 keys, got %d", len(res.Values))
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
}

func TestMergeMaps_StrategyLeft(t *testing.T) {
	left := map[string]string{"A": "left"}
	right := map[string]string{"A": "right"}
	res, err := MergeMaps(left, right, StrategyLeft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Values["A"] != "left" {
		t.Errorf("expected 'left', got %q", res.Values["A"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "A" {
		t.Errorf("expected conflict on A")
	}
}

func TestMergeMaps_StrategyRight(t *testing.T) {
	left := map[string]string{"A": "left"}
	right := map[string]string{"A": "right"}
	res, err := MergeMaps(left, right, StrategyRight)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Values["A"] != "right" {
		t.Errorf("expected 'right', got %q", res.Values["A"])
	}
}

func TestMergeMaps_StrategyError(t *testing.T) {
	left := map[string]string{"A": "left"}
	right := map[string]string{"A": "right"}
	_, err := MergeMaps(left, right, StrategyError)
	if err == nil {
		t.Fatal("expected error on conflict")
	}
}

func TestMergeMaps_SameValue_NoConflict(t *testing.T) {
	left := map[string]string{"A": "same"}
	right := map[string]string{"A": "same"}
	res, err := MergeMaps(left, right, StrategyError)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("same values should not be conflicts")
	}
}

func TestMergeMaps_EmptyInputs(t *testing.T) {
	res, err := MergeMaps(map[string]string{}, map[string]string{}, StrategyLeft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Values) != 0 {
		t.Errorf("expected empty result")
	}
}
