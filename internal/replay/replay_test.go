package replay_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/history"
	"github.com/user/stackdiff/internal/replay"
)

func makeRecord(t time.Time, status diff.Status) history.Record {
	return history.NewRecord("env-a", "env-b", []diff.Entry{
		{Key: "FOO", OldValue: "x", NewValue: "y", Status: status},
	}, history.WithTimestamp(t))
}

func buildStore(t *testing.T, records []history.Record) *history.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := history.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range records {
		if err := s.Append(r); err != nil {
			t.Fatal(err)
		}
	}
	return s
}

func TestRun_ReplayAll(t *testing.T) {
	now := time.Now()
	records := []history.Record{
		makeRecord(now.Add(-2*time.Hour), diff.StatusChanged),
		makeRecord(now.Add(-1*time.Hour), diff.StatusEqual),
	}
	store := buildStore(t, records)

	var count int
	err := replay.Run(store, func(r history.Record) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 records, got %d", count)
	}
}

func TestRun_OnlyDrift(t *testing.T) {
	now := time.Now()
	records := []history.Record{
		makeRecord(now.Add(-2*time.Hour), diff.StatusChanged),
		makeRecord(now.Add(-1*time.Hour), diff.StatusEqual),
	}
	store := buildStore(t, records)

	var count int
	err := replay.Run(store, func(r history.Record) error {
		count++
		return nil
	}, replay.WithOnlyDrift())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 drift record, got %d", count)
	}
}

func TestRun_SinceFilter(t *testing.T) {
	now := time.Now()
	records := []history.Record{
		makeRecord(now.Add(-3*time.Hour), diff.StatusChanged),
		makeRecord(now.Add(-1*time.Hour), diff.StatusChanged),
	}
	store := buildStore(t, records)

	var count int
	err := replay.Run(store, func(r history.Record) error {
		count++
		return nil
	}, replay.WithSince(now.Add(-2*time.Hour)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 record after since filter, got %d", count)
	}
}

func TestRun_HandlerError(t *testing.T) {
	now := time.Now()
	store := buildStore(t, []history.Record{makeRecord(now, diff.StatusChanged)})

	sentinel := errors.New("handler error")
	err := replay.Run(store, func(r history.Record) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestRun_NilStore(t *testing.T) {
	err := replay.Run(nil, func(r history.Record) error { return nil })
	if err == nil {
		t.Error("expected error for nil store")
	}
}

func TestRun_NilHandler(t *testing.T) {
	dir := t.TempDir()
	s, _ := history.NewStore(dir)
	err := replay.Run(s, nil)
	if err == nil {
		t.Error("expected error for nil handler")
	}
}
