// Package digest computes a deterministic hash fingerprint over a set of
// diff entries, making it easy to detect whether two runs produced identical
// drift output without comparing every field.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/yourorg/stackdiff/internal/diff"
)

// Result holds the computed fingerprint and entry count.
type Result struct {
	Hash       string
	EntryCount int
}

// String returns a short human-readable representation.
func (r Result) String() string {
	return fmt.Sprintf("digest:%s entries:%d", r.Hash[:12], r.EntryCount)
}

// Compute produces a stable SHA-256 fingerprint over the provided entries.
// Entries are sorted by key before hashing so order does not matter.
func Compute(entries []diff.Entry) Result {
	if len(entries) == 0 {
		return Result{Hash: emptyHash(), EntryCount: 0}
	}

	sorted := make([]diff.Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	h := sha256.New()
	for _, e := range sorted {
		fmt.Fprintf(h, "%s|%s|%s|%s\n", e.Key, e.OldVal, e.NewVal, e.Status)
	}

	return Result{
		Hash:       hex.EncodeToString(h.Sum(nil)),
		EntryCount: len(sorted),
	}
}

// Equal reports whether two Results represent identical drift sets.
func Equal(a, b Result) bool {
	return a.Hash == b.Hash
}

func emptyHash() string {
	h := sha256.New()
	return hex.EncodeToString(h.Sum(nil))
}
