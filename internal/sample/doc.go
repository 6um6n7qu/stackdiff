// Package sample provides rate-based and count-based sampling of diff entries.
//
// Use Apply to reduce a large set of drift entries to a representative subset,
// useful when feeding results into alerting or export pipelines that have
// volume constraints.
//
// Example:
//
//	cfg := sample.Config{Rate: 0.5, MaxEntries: 100, Seed: 42}
//	sampled := sample.Apply(entries, cfg)
package sample
