package compare_test

import (
	"testing"

	"github.com/yourorg/stackdiff/internal/compare"
	"github.com/yourorg/stackdiff/internal/diff"
)

func TestRun_NoDrift(t *testing.T) {
	a := map[string]string{"HOST": "localhost", "PORT": "8080"}
	b := map[string]string{"HOST": "localhost", "PORT": "8080"}

	r, err := compare.Run(a, b, compare.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestRun_DetectsDrift(t *testing.T) {
	a := map[string]string{"HOST": "localhost"}
	b := map[string]string{"HOST": "prod.example.com"}

	r, err := compare.Run(a, b, compare.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.HasDrift() {
		t.Error("expected drift")
	}
}

func TestRun_RedactsMasksSecrets(t *testing.T) {
	a := map[string]string{"DB_PASSWORD": "secret1"}
	b := map[string]string{"DB_PASSWORD": "secret2"}

	opts := compare.DefaultOptions()
	opts.RedactSecrets = true
	r, err := compare.Run(a, b, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range r.Entries {
		if e.OldVal == "secret1" || e.NewVal == "secret2" {
			t.Error("expected secrets to be masked")
		}
	}
}

func TestRun_FilterByStatus(t *testing.T) {
	a := map[string]string{"A": "1", "B": "old"}
	b := map[string]string{"A": "1", "B": "new", "C": "added"}

	opts := compare.DefaultOptions()
	opts.FilterStatuses = []string{string(diff.StatusChanged)}
	r, err := compare.Run(a, b, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range r.Entries {
		if e.Status != diff.StatusChanged {
			t.Errorf("expected only changed entries, got %s", e.Status)
		}
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := compare.DefaultOptions()
	if !opts.RedactSecrets {
		t.Error("expected RedactSecrets to be true by default")
	}
	if !opts.Normalize {
		t.Error("expected Normalize to be true by default")
	}
}
