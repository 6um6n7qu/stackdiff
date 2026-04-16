package watch

import (
	"context"
	"time"

	"github.com/user/stackdiff/internal/config"
	"github.com/user/stackdiff/internal/diff"
)

// Options configures the polling watcher.
type Options struct {
	Interval time.Duration
	FileA    string
	FileB    string
}

// Event is emitted whenever drift is detected between polls.
type Event struct {
	At      time.Time
	Entries []diff.Entry
}

// Watch polls two config files at the given interval and sends drift events
// on the returned channel. The channel is closed when ctx is cancelled.
func Watch(ctx context.Context, opts Options) (<-chan Event, error) {
	if opts.Interval <= 0 {
		opts.Interval = 30 * time.Second
	}

	events := make(chan Event, 4)

	go func() {
		defer close(events)
		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				entries, err := poll(opts.FileA, opts.FileB{
		\t	}
	ect {
				case events <- Event{At: time.Now(), Entries: entries}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return events, nil
}

func poll(fileA, fileB string) ([]diff.Entry, error) {
	cfgA, err := config.LoadFromFile(fileA)
	if err != nil {
		return nil, err
	}
	cfgB, err := config.LoadFromFile(fileB)
	if err != nil {
		return nil, err
	}
	entries := diff.Compare(cfgA, cfgB)
	var drifted []diff.Entry
	for _, e := range entries {
		if e.IsDrift() {
			drifted = append(drifted, e)
		}
	}
	return drifted, nil
}
