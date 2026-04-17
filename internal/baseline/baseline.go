// Package baseline provides functionality to capture and compare a known-good
// config state, allowing drift to be measured against a fixed reference point.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// Baseline represents a saved reference configuration state.
type Baseline struct {
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	Entries   []diff.Entry `json:"entries"`
}

// New creates a new Baseline with the given name and entries.
func New(name string, entries []diff.Entry) *Baseline {
	return &Baseline{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	}
}

// Save writes the baseline to a JSON file at the given path.
func Save(b *Baseline, path string) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// Load reads a baseline from a JSON file at the given path.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// DriftFrom compares the current entries against the baseline and returns
// any entries that differ from the known-good state.
func (b *Baseline) DriftFrom(current []diff.Entry) []diff.Entry {
	index := make(map[string]diff.Entry, len(b.Entries))
	for _, e := range b.Entries {
		index[e.Key] = e
	}
	var drifted []diff.Entry
	for _, c := range current {
		ref, ok := index[c.Key]
		if !ok || ref.NewValue != c.NewValue {
			drifted = append(drifted, c)
		}
	}
	return drifted
}
