package mask

import (
	"testing"

	"github.com/user/stackdiff/internal/diff"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "db_password", OldVal: "old", NewVal: "new", Status: diff.StatusChanged},
		{Key: "log_level", OldVal: "debug", NewVal: "info", Status: diff.StatusChanged},
		{Key: "api_token", OldVal: "", NewVal: "tok123", Status: diff.StatusAdded},
	}
}

func TestApplyToEntries_MasksSensitive(t *testing.T) {
	m := New()
	out := ApplyToEntries(m, sampleEntries())
	for _, e := range out {
		if m.IsSensitive(e.Key) {
			if e.OldVal != "" && e.OldVal != DefaultReplacement {
				t.Errorf("key %q OldVal not masked: %q", e.Key, e.OldVal)
			}
			if e.NewVal != "" && e.NewVal != DefaultReplacement {
				t.Errorf("key %q NewVal not masked: %q", e.Key, e.NewVal)
			}
		}
	}
}

func TestApplyToEntries_PreservesNonSensitive(t *testing.T) {
	m := New()
	out := ApplyToEntries(m, sampleEntries())
	for _, e := range out {
		if e.Key == "log_level" {
			if e.OldVal != "debug" || e.NewVal != "info" {
				t.Errorf("log_level values should not be masked")
			}
		}
	}
}

func TestApplyToEntries_EmptyOldValPreserved(t *testing.T) {
	m := New()
	out := ApplyToEntries(m, sampleEntries())
	for _, e := range out {
		if e.Key == "api_token" {
			if e.OldVal != "" {
				t.Errorf("empty OldVal should remain empty, got %q", e.OldVal)
			}
			if e.NewVal != DefaultReplacement {
				t.Errorf("non-empty NewVal should be masked, got %q", e.NewVal)
			}
		}
	}
}

func TestApplyToEntries_DoesNotMutateOriginal(t *testing.T) {
	m := New()
	orig := sampleEntries()
	ApplyToEntries(m, orig)
	if orig[0].NewVal != "new" {
		t.Error("original entries should not be mutated")
	}
}
