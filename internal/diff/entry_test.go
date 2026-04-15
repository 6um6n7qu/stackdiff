package diff_test

import (
	"testing"

	"github.com/yourusername/stackdiff/internal/diff"
)

func TestIsDrift_Equal(t *testing.T) {
	e := diff.Entry{Key: "PORT", Left: "8080", Right: "8080", Status: diff.StatusEqual}
	if e.IsDrift() {
		t.Error("expected IsDrift to be false for StatusEqual")
	}
}

func TestIsDrift_Changed(t *testing.T) {
	e := diff.Entry{Key: "ENV", Left: "dev", Right: "prod", Status: diff.StatusChanged}
	if !e.IsDrift() {
		t.Error("expected IsDrift to be true for StatusChanged")
	}
}

func TestIsDrift_Added(t *testing.T) {
	e := diff.Entry{Key: "NEW_KEY", Left: "", Right: "value", Status: diff.StatusAdded}
	if !e.IsDrift() {
		t.Error("expected IsDrift to be true for StatusAdded")
	}
}

func TestIsDrift_Removed(t *testing.T) {
	e := diff.Entry{Key: "OLD_KEY", Left: "value", Right: "", Status: diff.StatusRemoved}
	if !e.IsDrift() {
		t.Error("expected IsDrift to be true for StatusRemoved")
	}
}

func TestStatusConstants(t *testing.T) {
	statuses := []diff.Status{
		diff.StatusEqual,
		diff.StatusAdded,
		diff.StatusRemoved,
		diff.StatusChanged,
	}
	for _, s := range statuses {
		if string(s) == "" {
			t.Errorf("status constant should not be empty")
		}
	}
}
