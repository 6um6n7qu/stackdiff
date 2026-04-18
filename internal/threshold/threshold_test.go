package threshold_test

import (
	"testing"

	"github.com/yourusername/stackdiff/internal/score"
	"github.com/yourusername/stackdiff/internal/threshold"
)

func makeScore(s float64, grade string) score.Result {
	return score.Result{Score: s, Grade: grade}
}

func TestEvaluate_OK(t *testing.T) {
	cfg := threshold.DefaultConfig()
	r := threshold.Evaluate(makeScore(5.0, "A"), cfg)
	if r.Level != threshold.LevelOK {
		t.Errorf("expected OK, got %s", r.Level)
	}
	if r.Breached() {
		t.Error("expected Breached to be false")
	}
}

func TestEvaluate_Warning(t *testing.T) {
	cfg := threshold.DefaultConfig()
	r := threshold.Evaluate(makeScore(25.0, "C"), cfg)
	if r.Level != threshold.LevelWarning {
		t.Errorf("expected Warning, got %s", r.Level)
	}
	if !r.Breached() {
		t.Error("expected Breached to be true")
	}
}

func TestEvaluate_Critical(t *testing.T) {
	cfg := threshold.DefaultConfig()
	r := threshold.Evaluate(makeScore(60.0, "F"), cfg)
	if r.Level != threshold.LevelCritical {
		t.Errorf("expected Critical, got %s", r.Level)
	}
	if !r.Breached() {
		t.Error("expected Breached to be true")
	}
}

func TestEvaluate_ExactWarningBoundary(t *testing.T) {
	cfg := threshold.DefaultConfig()
	r := threshold.Evaluate(makeScore(20.0, "B"), cfg)
	if r.Level != threshold.LevelWarning {
		t.Errorf("expected Warning at boundary, got %s", r.Level)
	}
}

func TestEvaluate_ExactCriticalBoundary(t *testing.T) {
	cfg := threshold.DefaultConfig()
	r := threshold.Evaluate(makeScore(50.0, "D"), cfg)
	if r.Level != threshold.LevelCritical {
		t.Errorf("expected Critical at boundary, got %s", r.Level)
	}
}

func TestEvaluate_MessageContainsScore(t *testing.T) {
	cfg := threshold.DefaultConfig()
	r := threshold.Evaluate(makeScore(30.0, "C"), cfg)
	if r.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := threshold.DefaultConfig()
	if cfg.Warning >= cfg.Critical {
		t.Error("warning threshold should be less than critical")
	}
}
