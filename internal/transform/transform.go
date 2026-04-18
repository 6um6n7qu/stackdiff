// Package transform applies key/value transformations to diff entries.
package transform

import (
	"strings"

	"github.com/user/stackdiff/internal/diff"
)

// Func is a transformation applied to a single entry.
type Func func(e diff.Entry) diff.Entry

// Transformer holds an ordered list of transform functions.
type Transformer struct {
	fns []Func
}

// New returns a Transformer with the given functions applied in order.
func New(fns ...Func) *Transformer {
	return &Transformer{fns: fns}
}

// Apply runs all transform functions over each entry and returns the results.
func (t *Transformer) Apply(entries []diff.Entry) []diff.Entry {
	out := make([]diff.Entry, len(entries))
	for i, e := range entries {
		for _, fn := range t.fns {
			e = fn(e)
		}
		out[i] = e
	}
	return out
}

// LowercaseKeys returns a Func that lower-cases every entry key.
func LowercaseKeys() Func {
	return func(e diff.Entry) diff.Entry {
		e.Key = strings.ToLower(e.Key)
		return e
	}
}

// TrimSpace returns a Func that trims whitespace from OldVal and NewVal.
func TrimSpace() Func {
	return func(e diff.Entry) diff.Entry {
		e.OldVal = strings.TrimSpace(e.OldVal)
		e.NewVal = strings.TrimSpace(e.NewVal)
		return e
	}
}

// PrefixKey returns a Func that prepends a prefix to every key.
func PrefixKey(prefix string) Func {
	return func(e diff.Entry) diff.Entry {
		e.Key = prefix + e.Key
		return e
	}
}
