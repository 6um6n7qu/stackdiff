package fingerprint_test

import (
	"testing"

	"github.com/your-org/stackdiff/internal/diff"
	"github.com/your-org/stackdiff/internal/fingerprint"
)

func TestCompute_NonEmpty(t *testing.T) {
	m := map[string]string{"HOST": "localhost", "PORT": "8080"}
	r := fingerprint.Compute(m)
	if r.Hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if r.KeyCount != 2 {
		t.Errorf("expected KeyCount=2, got %d", r.KeyCount)
	}
}

func TestCompute_Empty(t *testing.T) {
	r := fingerprint.Compute(map[string]string{})
	if r.Hash == "" {
		t.Fatal("expected non-empty hash for empty map")
	}
	if r.KeyCount != 0 {
		t.Errorf("expected KeyCount=0, got %d", r.KeyCount)
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	a := fingerprint.Compute(map[string]string{"A": "1", "B": "2"})
	b := fingerprint.Compute(map[string]string{"B": "2", "A": "1"})
	if !fingerprint.Equal(a, b) {
		t.Errorf("expected equal fingerprints, got %s vs %s", a, b)
	}
}

func TestCompute_DifferentValues(t *testing.T) {
	a := fingerprint.Compute(map[string]string{"HOST": "localhost"})
	b := fingerprint.Compute(map[string]string{"HOST": "prod.example.com"})
	if fingerprint.Equal(a, b) {
		t.Error("expected different fingerprints for different values")
	}
}

func TestFromEntries_UsesNewValue(t *testing.T) {
	entries := []diff.Entry{
		{Key: "HOST", OldValue: "old", NewValue: "new", Status: diff.StatusChanged},
		{Key: "PORT", OldValue: "8080", NewValue: "8080", Status: diff.StatusEqual},
	}
	r := fingerprint.FromEntries(entries)
	if r.Hash == "" {
		t.Fatal("expected non-empty hash from entries")
	}
	if r.KeyCount != 2 {
		t.Errorf("expected KeyCount=2, got %d", r.KeyCount)
	}
}

func TestEqual_SameInput(t *testing.T) {
	m := map[string]string{"X": "y"}
	a := fingerprint.Compute(m)
	b := fingerprint.Compute(m)
	if !fingerprint.Equal(a, b) {
		t.Error("expected equal fingerprints for same input")
	}
}

func TestResult_String(t *testing.T) {
	r := fingerprint.Compute(map[string]string{"K": "V"})
	if r.String() != r.Hash {
		t.Error("String() should return Hash")
	}
}
