package throttle

import (
	"fmt"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// Result holds the outcome of a throttle pipeline run.
type Result struct {
	Allowed  bool
	Entries  []diff.Entry
	Dropped  int
	Message  string
}

// Run applies throttle logic to a slice of drift entries.
// If the event is allowed, entries are returned unchanged.
// If throttled, Allowed is false and Entries is nil.
func Run(t *Throttle, entries []diff.Entry) Result {
	now := time.Now()
	if !t.Allow(now) {
		return Result{
			Allowed: false,
			Entries: nil,
			Dropped: len(entries),
			Message: fmt.Sprintf("throttled: exceeded %d events in %s", t.cfg.MaxEvents, t.cfg.Window),
		}
	}
	return Result{
		Allowed: true,
		Entries: entries,
		Dropped: 0,
	}
}

// MustAllow panics if the throttle gate is closed, otherwise returns entries.
func MustAllow(t *Throttle, entries []diff.Entry) []diff.Entry {
	res := Run(t, entries)
	if !res.Allowed {
		panic(res.Message)
	}
	return res.Entries
}
