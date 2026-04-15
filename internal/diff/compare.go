package diff

// Kind represents the type of difference between two configs.
type Kind string

const (
	Added   Kind = "added"
	Removed Kind = "removed"
	Changed Kind = "changed"
)

// DiffEntry describes a single key-level difference.
type DiffEntry struct {
	Key      string `json:"key"`
	Kind     Kind   `json:"kind"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

// Compare returns the list of differences between source and target config maps.
// Keys present only in target are Added; only in source are Removed; in both
// but with different values are Changed.
func Compare(source, target map[string]string) []DiffEntry {
	var entries []DiffEntry

	for k, tv := range target {
		if sv, ok := source[k]; !ok {
			entries = append(entries, DiffEntry{Key: k, Kind: Added, NewValue: tv})
		} else if sv != tv {
			entries = append(entries, DiffEntry{Key: k, Kind: Changed, OldValue: sv, NewValue: tv})
		}
	}

	for k, sv := range source {
		if _, ok := target[k]; !ok {
			entries = append(entries, DiffEntry{Key: k, Kind: Removed, OldValue: sv})
		}
	}

	sortEntries(entries)
	return entries
}

func sortEntries(entries []DiffEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].Key < entries[j-1].Key; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}
