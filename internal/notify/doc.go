// Package notify provides drift notification support for stackdiff.
//
// Supported channels:
//   - stdout  — writes a human-readable summary to standard output (default)
//   - slack   — posts a message to a Slack incoming webhook URL
//   - file    — reserved for future file-based notification
//
// Usage:
//
//	n := notify.New(notify.DefaultConfig())
//	if err := n.Notify(entries); err != nil {
//		log.Fatal(err)
//	}
package notify
