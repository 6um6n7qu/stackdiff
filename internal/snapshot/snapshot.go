package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/stackdiff/internal/diff"
)

// Snapshot captures a diff result at a point in time.
type Snapshot struct {
	Timestamp time.Time    `json:"timestamp"`
	LabelA    string       `json:"label_a"`
	LabelB    string       `json:"label_b"`
	Entries   []diff.Entry `json:"entries"`
}

// New creates a new Snapshot from the provided diff entries and labels.
func New(labelA, labelB string, entries []diff.Entry) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		LabelA:    labelA,
		LabelB:    labelB,
		Entries:   entries,
	}
}

// Save writes the snapshot as JSON to the given file path.
func Save(s *Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file %q: %w", path, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}
