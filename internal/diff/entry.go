package diff

// Status represents the kind of difference found for a config key.
type Status string

const (
	// StatusEqual means both sides have the same value.
	StatusEqual Status = "equal"
	// StatusAdded means the key exists only on the right side.
	StatusAdded Status = "added"
	// StatusRemoved means the key exists only on the left side.
	StatusRemoved Status = "removed"
	// StatusChanged means the key exists on both sides but values differ.
	StatusChanged Status = "changed"
)

// Entry represents a single compared config key with its diff result.
type Entry struct {
	Key    string
	Left   string
	Right  string
	Status Status
}

// IsDrift returns true if the entry represents any form of configuration drift.
func (e Entry) IsDrift() bool {
	return e.Status != StatusEqual
}
