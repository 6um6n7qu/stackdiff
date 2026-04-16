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
package watch
