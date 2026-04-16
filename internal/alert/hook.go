package alert

import "github.com/stackdiff/internal/diff"

// HookFunc is a callback invoked when drift is detected.
type HookFunc func(a *Alert)

// Dispatcher holds a set of hooks and fires them on drift.
type Dispatcher struct {
	hooks []HookFunc
	cfg   Config
}

// NewDispatcher creates a Dispatcher with the given config.
func NewDispatcher(cfg Config) *Dispatcher {
	return &Dispatcher{cfg: cfg}
}

// Register adds a hook to the dispatcher.
func (d *Dispatcher) Register(fn HookFunc) {
	d.hooks = append(d.hooks, fn)
}

// Dispatch emits an alert for entries and calls all registered hooks.
func (d *Dispatcher) Dispatch(entries []diff.Entry) *Alert {
	a := Emit(entries, d.cfg)
	if a == nil {
		return nil
	}
	for _, fn := range d.hooks {
		fn(a)
	}
	return a
}
