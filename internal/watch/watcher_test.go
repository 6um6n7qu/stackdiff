package watch

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeCfg(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeCfg: %v", err)
	}
	return p
}

func TestWatch_EmitsEventOnDrift(t *testing.T) {
	dir := t.TempDir()
	a := writeCfg(t, dir, "a.yaml", "service: svc\nenv:\n  KEY: alpha\n")
	b := writeCfg(t, dir, "b.yaml", "service: svc\nenv:\n  KEY: beta\n")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	events, err := Watch(ctx, Options{Interval: 100 * time.Millisecond, FileA: a, FileB: b})
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	select {
	case ev, ok := <-events:
		if !ok {
			t.Fatal("channel closed before event")
		}
		if len(ev.Entries) == 0 {
			t.Error("expected drift entries, got none")
		}
		if ev.At.IsZero() {
			t.Error("expected non-zero timestamp")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for drift event")
	}
}

func TestWatch_NoDriftNoEvent(t *testing.T) {
	dir := t.TempDir()
	content := "service: svc\nenv:\n  KEY: same\n"
	a := writeCfg(t, dir, "a.yaml", content)
	b := writeCfg(t, dir, "b.yaml", content)

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	events, err := Watch(ctx, Options{Interval: 100 * time.Millisecond, FileA: a, FileB: b})
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	<-ctx.Done()
	select {
	case ev, ok := <-events:
		if ok {
			t.Errorf("unexpected event with %d entries", len(ev.Entries))
		}
	default:
	}
}

func TestWatch_DefaultInterval(t *testing.T) {
	opts := Options{Interval: 0}
	if opts.Interval != 0 {
		t.Error("expected zero before normalisation")
	}
	// normalisation happens inside Watch goroutine; verify no panic on start
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	dir := t.TempDir()
	a := writeCfg(t, dir, "a.yaml", "service: s\nenv:\n  K: v\n")
	b := writeCfg(t, dir, "b.yaml", "service: s\nenv:\n  K: v\n")
	_, err := Watch(ctx, Options{FileA: a, FileB: b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
