// Package emit provides a unified event emission layer for stackdiff.
// It allows components to publish drift events to multiple registered
// sinks (e.g. stdout, webhook, file) through a single Emitter interface.
package emit

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/stackdiff/stackdiff/internal/diff"
)

// Event represents a single drift emission with metadata.
type Event struct {
	// ID is a short identifier for the event (e.g. a run ID or label).
	ID string

	// OccurredAt is when the drift was detected.
	OccurredAt time.Time

	// Entries holds the diff entries associated with this event.
	Entries []diff.Entry

	// Labels are arbitrary key/value metadata attached to the event.
	Labels map[string]string
}

// HasDrift returns true if any entry in the event represents a drift.
func (e Event) HasDrift() bool {
	for _, en := range e.Entries {
		if en.IsDrift() {
			return true
		}
	}
	return false
}

// Sink is the interface implemented by anything that can receive an Event.
type Sink interface {
	Receive(Event) error
}

// WriterSink writes a human-readable summary of the event to an io.Writer.
type WriterSink struct {
	w io.Writer
}

// NewWriterSink returns a Sink that writes events to w.
func NewWriterSink(w io.Writer) *WriterSink {
	return &WriterSink{w: w}
}

// Receive writes a one-line summary of the event to the underlying writer.
func (s *WriterSink) Receive(ev Event) error {
	driftCount := 0
	for _, e := range ev.Entries {
		if e.IsDrift() {
			driftCount++
		}
	}
	_, err := fmt.Fprintf(s.w, "[%s] event=%s drift=%d total=%d\n",
		ev.OccurredAt.Format(time.RFC3339),
		ev.ID,
		driftCount,
		len(ev.Entries),
	)
	return err
}

// Emitter fans out events to all registered sinks.
type Emitter struct {
	mu    sync.RWMutex
	sinks []Sink
}

// New returns a new Emitter with no sinks registered.
func New() *Emitter {
	return &Emitter{}
}

// Register adds a Sink to the Emitter. It is safe to call concurrently.
func (em *Emitter) Register(s Sink) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.sinks = append(em.sinks, s)
}

// Emit sends ev to all registered sinks. Errors from individual sinks are
// collected and returned as a combined error; remaining sinks still receive
// the event even if an earlier sink fails.
func (em *Emitter) Emit(ev Event) error {
	em.mu.RLock()
	sinks := make([]Sink, len(em.sinks))
	copy(sinks, em.sinks)
	em.mu.RUnlock()

	if ev.OccurredAt.IsZero() {
		ev.OccurredAt = time.Now()
	}

	var errs []error
	for _, s := range sinks {
		if err := s.Receive(ev); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("emit: %d sink(s) failed: %v", len(errs), errs)
}

// EmitIfDrift is a convenience wrapper that only calls Emit when the event
// contains at least one drifted entry.
func (em *Emitter) EmitIfDrift(ev Event) error {
	if !ev.HasDrift() {
		return nil
	}
	return em.Emit(ev)
}
