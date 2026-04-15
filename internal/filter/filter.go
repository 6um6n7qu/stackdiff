package filter

import (
	"strings"

	"github.com/yourusername/stackdiff/internal/diff"
)

// Options holds filtering criteria for diff entries.
type Options struct {
	OnlyChanged bool
	OnlyAdded   bool
	OnlyRemoved bool
	KeyPrefix   string
}

// Apply filters a slice of diff.Entry values based on the provided Options.
// If no filter flags are set, all entries are returned unchanged.
func Apply(entries []diff.Entry, opts Options) []diff.Entry {
	var result []diff.Entry

	for _, e := range entries {
		if opts.KeyPrefix != "" && !strings.HasPrefix(e.Key, opts.KeyPrefix) {
			continue
		}

		if !opts.OnlyChanged && !opts.OnlyAdded && !opts.OnlyRemoved {
			result = append(result, e)
			continue
		}

		switch e.Status {
		case diff.StatusChanged:
			if opts.OnlyChanged {
				result = append(result, e)
			}
		case diff.StatusAdded:
			if opts.OnlyAdded {
				result = append(result, e)
			}
		case diff.StatusRemoved:
			if opts.OnlyRemoved {
				result = append(result, e)
			}
		}
	}

	return result
}
