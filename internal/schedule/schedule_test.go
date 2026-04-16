package schedule_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/schedule"
)

func TestRun_CallsOnDriftWhenDriftDetected(t *testing.T) {
	var mu sync.Mutex
	var got []diff.Entry

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	job := schedule.Job{
		Interval: 50 * time.Millisecond,
		Loader: func() (map[string]string, map[string]string, error) {
			return map[string]string{"KEY": "a"}, map[string]string{"KEY": "b"}, nil
		},
		OnDrift: func(entries []diff.Entry) {
			mu.Lock()
			got = entries
			mu.Unlock()
			cancel()
		},
	}
	schedule.Run(ctx, job)

	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatal("expected drift entries, got none")
	}
}

func TestRun_NoDriftNoCallback(t *testing.T) {
	called := false
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	job := schedule.Job{
		Interval: 40 * time.Millisecond,
		Loader: func() (map[string]string, map[string]string, error) {
			return map[string]string{"K": "v"}, map[string]string{"K": "v"}, nil
		},
		OnDrift: func(_ []diff.Entry) { called = true },
	}
	schedule.Run(ctx, job)
	if called {
		t.Fatal("OnDrift should not be called when there is no drift")
	}
}

func TestRun_CallsOnErrorOnLoaderFailure(t *testing.T) {
	var errGot error
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	job := schedule.Job{
		Interval: 40 * time.Millisecond,
		Loader: func() (map[string]string, map[string]string, error) {
			return nil, nil, errors.New("load failed")
		},
		OnError: func(err error) {
			errGot = err
			cancel()
		},
	}
	schedule.Run(ctx, job)
	if errGot == nil {
		t.Fatal("expected error callback to be invoked")
	}
}

func TestNewJob_DefaultInterval(t *testing.T) {
	j := schedule.NewJob(nil)
	if j.Interval != time.Minute {
		t.Fatalf("expected 1m, got %v", j.Interval)
	}
}

func TestWithInterval_Override(t *testing.T) {
	j := schedule.NewJob(nil, schedule.WithInterval(5*time.Second))
	if j.Interval != 5*time.Second {
		t.Fatalf("expected 5s, got %v", j.Interval)
	}
}
