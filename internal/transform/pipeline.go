package transform

import "github.com/user/stackdiff/internal/diff"

// Pipeline applies a sequence of named transformers to a slice of entries.
type Pipeline struct {
	steps []*Transformer
}

// NewPipeline constructs a Pipeline from the provided transformers.
func NewPipeline(steps ...*Transformer) *Pipeline {
	return &Pipeline{steps: steps}
}

// Run executes every transformer in order and returns the final entries.
func (p *Pipeline) Run(entries []diff.Entry) []diff.Entry {
	result := entries
	for _, t := range p.steps {
		result = t.Apply(result)
	}
	return result
}

// DefaultPipeline returns a Pipeline with sensible defaults:
// trim whitespace then lowercase keys.
func DefaultPipeline() *Pipeline {
	return NewPipeline(
		New(TrimSpace()),
		New(LowercaseKeys()),
	)
}
