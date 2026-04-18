// Package ignore provides key-based filtering to suppress known drift entries.
package ignore

import (
	"strings"

	"github.com/user/stackdiff/internal/diff"
)

// List holds a set of keys and key prefixes to ignore during diff evaluation.
type List struct {
	keys    map[string]struct{}
	prefixes []string
}

// New creates a new List from exact keys and prefix patterns.
// Prefix patterns should end with "*" (e.g. "DEBUG_*").
func New(patterns []string) *List {
	l := &List{keys: make(map[string]struct{})}
	for _, p := range patterns {
		if strings.HasSuffix(p, "*") {
			l.prefixes = append(l.prefixes, strings.TrimSuffix(p, "*"))
		} else {
			l.keys[p] = struct{}{}
		}
	}
	return l
}

// Match reports whether the given key should be ignored.
func (l *List) Match(key string) bool {
	if _, ok := l.keys[key]; ok {
		return true
	}
	for _, prefix := range l.prefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// Apply removes entries whose keys match the ignore list.
func (l *List) Apply(entries []diff.Entry) []diff.Entry {
	out := make([]diff.Entry, 0, len(entries))
	for _, e := range entries {
		if !l.Match(e.Key) {
			out = append(out, e)
		}
	}
	return out
}
