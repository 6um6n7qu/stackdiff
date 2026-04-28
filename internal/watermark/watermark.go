// Package watermark tracks high-water marks for drift counts across runs,
// allowing callers to detect when drift volume has reached a new peak.
package watermark

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// Mark holds the recorded peak drift count and when it was observed.
type Mark struct {
	Peak      int       `json:"peak"`
	ObservedAt time.Time `json:"observed_at"`
}

// Store persists and retrieves watermark data for a named series.
type Store struct {
	mu   sync.Mutex
	path string
	data map[string]Mark
}

// New creates a Store backed by the given file path.
func New(path string) *Store {
	return &Store{path: path, data: make(map[string]Mark)}
}

// Load reads persisted marks from disk. Missing file is not an error.
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.data)
}

// Save writes current marks to disk.
func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}

// Record updates the high-water mark for series if count exceeds the current peak.
// Returns true when a new peak is set.
func (s *Store) Record(series string, count int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.data[series]
	if !ok || count > current.Peak {
		s.data[series] = Mark{Peak: count, ObservedAt: time.Now().UTC()}
		return true
	}
	return false
}

// Get returns the current high-water mark for a series.
func (s *Store) Get(series string) (Mark, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	m, ok := s.data[series]
	return m, ok
}

// Reset clears the mark for a series.
func (s *Store) Reset(series string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, series)
}
