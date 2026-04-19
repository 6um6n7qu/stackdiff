package diff

// Chain applies a sequence of named transformations to a slice of entries,
// returning the result after each step. Useful for pipeline introspection.

// Step represents a named transformation applied to entries.
type Step struct {
	Name    string
	Apply   func([]Entry) []Entry
}

// Chain runs entries through each Step in order and returns the final result.
func Chain(entries []Entry, steps []Step) []Entry {
	current := entries
	for _, s := range steps {
		current = s.Apply(current)
	}
	return current
}

// ChainWithTrace runs entries through each Step and returns intermediate results
// keyed by step name.
func ChainWithTrace(entries []Entry, steps []Step) map[string][]Entry {
	trace := make(map[string][]Entry, len(steps))
	current := entries
	for _, s := range steps {
		current = s.Apply(current)
		copy := make([]Entry, len(current))
		for i, e := range current {
			copy[i] = e
		}
		trace[s.Name] = copy
	}
	return trace
}
