// Package threshold provides drift alerting based on configurable score limits.
package threshold

import (
	"fmt"

	"github.com/yourusername/stackdiff/internal/score"
)

// Level represents the severity of a threshold breach.
type Level string

const (
	LevelOK      Level = "ok"
	LevelWarning Level = "warning"
	LevelCritical Level = "critical"
)

// Config holds warning and critical score thresholds.
type Config struct {
	Warning  float64
	Critical float64
}

// DefaultConfig returns sensible default thresholds.
func DefaultConfig() Config {
	return Config{
		Warning:  20.0,
		Critical: 50.0,
	}
}

// Result holds the evaluation outcome.
type Result struct {
	Level   Level
	Score   float64
	Message string
}

// Evaluate checks a score.Result against the configured thresholds.
func Evaluate(s score.Result, cfg Config) Result {
	switch {
	case s.Score >= cfg.Critical:
		return Result{
			Level:   LevelCritical,
			Score:   s.Score,
			Message: fmt.Sprintf("critical drift detected: score %.1f (grade %s)", s.Score, s.Grade),
		}
	case s.Score >= cfg.Warning:
		return Result{
			Level:   LevelWarning,
			Score:   s.Score,
			Message: fmt.Sprintf("warning: elevated drift score %.1f (grade %s)", s.Score, s.Grade),
		}
	default:
		return Result{
			Level:   LevelOK,
			Score:   s.Score,
			Message: fmt.Sprintf("drift within acceptable range: score %.1f (grade %s)", s.Score, s.Grade),
		}
	}
}

// Breached returns true if the level is warning or critical.
func (r Result) Breached() bool {
	return r.Level != LevelOK
}
