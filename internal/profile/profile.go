package profile

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Profile represents a named environment profile with metadata.
type Profile struct {
	Name    string            `yaml:"name"`
	Desc    string            `yaml:"desc"`
	EnvFile string            `yaml:"env_file"`
	Labels  map[string]string `yaml:"labels"`
}

// Store manages named profiles on disk.
type Store struct {
	Dir string
}

// NewStore returns a Store rooted at dir.
func NewStore(dir string) *Store {
	return &Store{Dir: dir}
}

// Save writes a profile to disk as <name>.yaml.
func (s *Store) Save(p *Profile) error {
	if err := os.MkdirAll(s.Dir, 0755); err != nil {
		return fmt.Errorf("profile: mkdir: %w", err)
	}
	path := filepath.Join(s.Dir, p.Name+".yaml")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("profile: create: %w", err)
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(p)
}

// Load reads a profile by name from disk.
func (s *Store) Load(name string) (*Profile, error) {
	path := filepath.Join(s.Dir, name+".yaml")
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("profile: open: %w", err)
	}
	defer f.Close()
	var p Profile
	if err := yaml.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("profile: decode: %w", err)
	}
	return &p, nil
}

// List returns all profile names available in the store.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("profile: readdir: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".yaml" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
