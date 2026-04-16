// Package watch provides a polling watcher that periodically compares two
// config files and emits drift events when differences are detected.
//
// Basic usage:
//
//	events, err := watch.Watch(ctx, watch.Options{
//		Interval: 10 * time.Second,
//		FileA:    "prod.yaml",
//		FileB:    "staging.yaml",
//	})
//	for e := range events {
//		fmt.Printf("drift detected at %s: %d keys\n", e.At, len(e.Entries))
//	}
//
// The returned channel is closed when the provided context is cancelled or
// when a fatal read error occurs. Callers should check ctx.Err() after the
// channel closes to distinguish between normal shutdown and an error.
//
// FileA and FileB may be the same path; in that case no drift will ever be
// reported, which can be useful in tests.
package watch
