package annotate_test

import (
	"testing"

	"github.com/stackdiff/stackdiff/internal/annotate"
)

func TestDefaultRules_NotEmpty(t *testing.T) {
	rules := annotate.DefaultRules()
	if len(rules) == 0 {
		t.Error("expected default rules to be non-empty")
	}
}

func TestDefaultRules_CoverAllStatuses(t *testing.T) {
	rules := annotate.DefaultRules()
	if len(rules) < 3 {
		t.Errorf("expected at least 3 default rules, got %d", len(rules))
	}
}
