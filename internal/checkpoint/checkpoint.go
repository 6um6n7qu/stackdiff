// Package checkpoint provides a mechanism for marking and comparing
// named points in time within a drift detection workflow. A checkpoint
// captures a labelled snapshot of config entries so that future runs
// can compare against a known-good state rather than the live config.
package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// Checkpoint holds a named, timestamped set of config entries.
type Checkpoint struct {
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	Entries   []diff.Entry `json:"entries"`
}

// New creates a new Checkpoint with the given name and entries.
func New(name string, entries []diff.Entry) *Checkpoint {
	return &Checkpoint{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	}
}

// Save writes the checkpoint to a JSON file under dir/<name>.json.
func Save(dir, name string, cp *Checkpoint) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("checkpoint: mkdir %s: %w", dir, err)
	}
	path := filepath.Join(dir, name+".json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("checkpoint: create %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cp); err != nil {
		return fmt.Errorf("checkpoint: encode: %w", err)
	}
	return nil
}

// Load reads a checkpoint from dir/<name>.json.
func Load(dir, name string) (*Checkpoint, error) {
	path := filepath.Join(dir, name+".json")
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("checkpoint: open %s: %w", path, err)
	}
	defer f.Close()
	var cp Checkpoint
	if err := json.NewDecoder(f).Decode(&cp); err != nil {
		return nil, fmt.Errorf("checkpoint: decode: %w", err)
	}
	return &cp, nil
}

// List returns the names of all checkpoints stored under dir.
func List(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("checkpoint: readdir %s: %w", dir, err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
