// Package history tracks past diff runs and allows retrieval of previous comparisons.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

// Record represents a single historical diff run.
type Record struct {
	ID        string       `json:"id"`
	Timestamp time.Time    `json:"timestamp"`
	LeftLabel string       `json:"left_label"`
	RightLabel string      `json:"right_label"`
	Entries   []diff.Entry `json:"entries"`
}

// Store manages persistence of history records on disk.
type Store struct {
	Dir string
}

// NewStore creates a Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("history: create dir: %w", err)
	}
	return &Store{Dir: dir}, nil
}

// Save persists a Record to disk as a JSON file named by its ID.
func (s *Store) Save(r Record) error {
	path := filepath.Join(s.Dir, r.ID+".json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("history: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("history: encode record: %w", err)
	}
	return nil
}

// Load retrieves a Record by ID from the store.
func (s *Store) Load(id string) (Record, error) {
	path := filepath.Join(s.Dir, id+".json")
	f, err := os.Open(path)
	if err != nil {
		return Record{}, fmt.Errorf("history: open file: %w", err)
	}
	defer f.Close()
	var r Record
	if err := json.NewDecoder(f).Decode(&r); err != nil {
		return Record{}, fmt.Errorf("history: decode record: %w", err)
	}
	return r, nil
}

// List returns all Record IDs stored in the directory, sorted by filename.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		return nil, fmt.Errorf("history: read dir: %w", err)
	}
	var ids []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			ids = append(ids, e.Name()[:len(e.Name())-5])
		}
	}
	return ids, nil
}
