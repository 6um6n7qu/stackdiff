package diff

import (
	"testing"

	"github.com/stackdiff/stackdiff/internal/config"
)

func TestCompare_NoDrift(t *testing.T) {
	a := config.Config{"KEY": "value", "PORT": "8080"}
	b := config.Config{"KEY": "value", "PORT": "8080"}

	result := Compare(a, b)
	if len(result) != 0 {
		t.Errorf("expected 0 drift entries, got %d", len(result))
	}
}

func TestCompare_Added(t *testing.T) {
	a := config.Config{"KEY": "value"}
	b := config.Config{"KEY": "value", "NEW_KEY": "new"}

	result := Compare(a, b)
	if len(result) != 1 {
		t.Fatalf("expected 1 drift entry, got %d", len(result))
	}
	if result[0].Kind != Added || result[0].Key != "NEW_KEY" {
		t.Errorf("unexpected entry: %+v", result[0])
	}
}

func TestCompare_Removed(t *testing.T) {
	a := config.Config{"KEY": "value", "OLD_KEY": "old"}
	b := config.Config{"KEY": "value"}

	result := Compare(a, b)
	if len(result) != 1 {
		t.Fatalf("expected 1 drift entry, got %d", len(result))
	}
	if result[0].Kind != Removed || result[0].Key != "OLD_KEY" {
		t.Errorf("unexpected entry: %+v", result[0])
	}
}

func TestCompare_Changed(t *testing.T) {
	a := config.Config{"DB_HOST": "localhost"}
	b := config.Config{"DB_HOST": "prod.db.internal"}

	result := Compare(a, b)
	if len(result) != 1 {
		t.Fatalf("expected 1 drift entry, got %d", len(result))
	}
	e := result[0]
	if e.Kind != Changed || e.ValueA != "localhost" || e.ValueB != "prod.db.internal" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestCompare_Mixed(t *testing.T) {
	a := config.Config{"A": "1", "B": "2", "C": "3"}
	b := config.Config{"A": "1", "B": "changed", "D": "4"}

	result := Compare(a, b)
	if len(result) != 3 {
		t.Errorf("expected 3 drift entries, got %d", len(result))
	}
}
