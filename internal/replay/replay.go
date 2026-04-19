// Package replay provides utilities for replaying historical diff records
// through a processing pipeline for audit, analysis, or re-export purposes.
package replay

import (
	"errors"
	"time"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/history"
)

// Options controls replay behaviour.
type Options struct {
	Since  time.Time
	Until  time.Time
	OnlyDrift bool
}

// Handler is called for each replayed record.
type Handler func(r history.Record) error

// Option is a functional option for Options.
type Option func(*Options)

// WithSince filters records after t.
func WithSince(t time.Time) Option {
	return func(o *Options) { o.Since = t }
}

// WithUntil filters records before t.
func WithUntil(t time.Time) Option {
	return func(o *Options) { o.Until = t }
}

// WithOnlyDrift skips records with no drift.
func WithOnlyDrift() Option {
	return func(o *Options) { o.OnlyDrift = true }
}

// Run replays records from store through handler using the given options.
func Run(store *history.Store, handler Handler, opts ...Option) error {
	if store == nil {
		return errors.New("replay: store must not be nil")
	}
	if handler == nil {
		return errors.New("replay: handler must not be nil")
	}

	cfg := &Options{}
	for _, o := range opts {
		o(cfg)
	}

	records, err := store.All()
	if err != nil {
		return err
	}

	for _, r := range records {
		if !cfg.Since.IsZero() && r.Timestamp.Before(cfg.Since) {
			continue
		}
		if !cfg.Until.IsZero() && r.Timestamp.After(cfg.Until) {
			continue
		}
		if cfg.OnlyDrift && !hasDrift(r.Entries) {
			continue
		}
		if err := handler(r); err != nil {
			return err
		}
	}
	return nil
}

func hasDrift(entries []diff.Entry) bool {
	for _, e := range entries {
		if e.IsDrift() {
			return true
		}
	}
	return false
}
