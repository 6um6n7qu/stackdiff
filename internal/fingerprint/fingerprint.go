// Package fingerprint generates a stable identity string for a config state,
// useful for detecting whether two snapshots represent the same effective config.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/your-org/stackdiff/internal/diff"
)

// Result holds the computed fingerprint and metadata.
type Result struct {
	Hash    string
	KeyCount int
}

// String returns the hex fingerprint hash.
func (r Result) String() string {
	return r.Hash
}

// Equal reports whether two fingerprints represent the same config state.
func Equal(a, b Result) bool {
	return a.Hash == b.Hash
}

// Compute derives a fingerprint from a map of config key-value pairs.
// Keys are sorted before hashing to ensure order-independence.
func Compute(entries map[string]string) Result {
	if len(entries) == 0 {
		return Result{Hash: emptyHash(), KeyCount: 0}
	}

	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, entries[k])
	}

	return Result{
		Hash:    hex.EncodeToString(h.Sum(nil)),
		KeyCount: len(keys),
	}
}

// FromEntries derives a fingerprint from a slice of diff.Entry values,
// using only entries that are not drift (i.e., status Equal).
func FromEntries(entries []diff.Entry) Result {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.NewValue
	}
	return Compute(m)
}

func emptyHash() string {
	h := sha256.New()
	return hex.EncodeToString(h.Sum(nil))
}
