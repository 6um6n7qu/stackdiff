package compare_test

import (
	"testing"

	"github.com/yourorg/stackdiff/internal/compare"
)

func TestWithRedact(t *testing.T) {
	opts := compare.WithRedact(compare.DefaultOptions(), false)
	if opts.RedactSecrets {
		t.Error("expected RedactSecrets false")
	}
}

func TestWithNormalize(t *testing.T) {
	opts := compare.WithNormalize(compare.DefaultOptions(), false)
	if opts.Normalize {
		t.Error("expected Normalize false")
	}
}

func TestWithKeyPrefix(t *testing.T) {
	opts := compare.WithKeyPrefix(compare.DefaultOptions(), "APP_")
	if opts.KeyPrefix != "APP_" {
		t.Errorf("expected APP_, got %s", opts.KeyPrefix)
	}
}

func TestWithFilterStatuses(t *testing.T) {
	opts := compare.WithFilterStatuses(compare.DefaultOptions(), "added", "removed")
	if len(opts.FilterStatuses) != 2 {
		t.Errorf("expected 2 statuses, got %d", len(opts.FilterStatuses))
	}
}
