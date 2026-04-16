package schedule_test

import (
	"testing"
	"time"

	"github.com/user/stackdiff/internal/schedule"
)

func TestWithOnError_SetsCallback(t *testing.T) {
	var captured error
	j := schedule.NewJob(nil, schedule.WithOnError(func(err error) {
		captured = err
	}))
	if j.OnError == nil {
		t.Fatal("expected OnError to be set")
	}
	sentinel := &struct{ error }{}
	_ = sentinel
	// just confirm it's callable
	j.OnError(nil)
	_ = captured
}

func TestNewJob_AppliesMultipleOptions(t *testing.T) {
	j := schedule.NewJob(
		nil,
		schedule.WithInterval(10*time.Second),
		schedule.WithOnError(func(error) {}),
	)
	if j.Interval != 10*time.Second {
		t.Fatalf("unexpected interval: %v", j.Interval)
	}
	if j.OnError == nil {
		t.Fatal("OnError should be set")
	}
}
