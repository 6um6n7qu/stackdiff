// Package tag assigns metadata labels to diff entries based on key prefix rules.
//
// Rules map key prefixes to human-readable tags (e.g. "db_" → "database").
// Tags can be used downstream for filtering, reporting, or alerting.
//
// Example:
//
//	tagger := tag.New(tag.DefaultRules())
//	labels := tagger.Apply(entries)
package tag
