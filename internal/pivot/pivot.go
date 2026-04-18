// Package pivot provides utilities for transposing diff entries
// into a key-centric view across multiple named environments.
package pivot

import (
	"sort"

	"github.com/user/stackdiff/internal/diff"
)

// Row represents a single key's values across multiple environments.
type Row struct {
	Key    string
	Values map[string]string // env name -> value
	Drift  bool
}

// Table is the result of pivoting entries from several environments.
type Table struct {
	Envs []string
	Rows []Row
}

// Build constructs a Table from a map of environment name to diff entries.
// Entries with StatusEqual are included to provide full context.
func Build(envEntries map[string][]diff.Entry) Table {
	keys := map[string]struct{}{}
	for _, entries := range envEntries {
		for _, e := range entries {
			keys[e.Key] = struct{}{}
		}
	}

	envNames := make([]string, 0, len(envEntries))
	for name := range envEntries {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	// index entries per env
	indexed := map[string]map[string]diff.Entry{}
	for name, entries := range envEntries {
		m := make(map[string]diff.Entry, len(entries))
		for _, e := range entries {
			m[e.Key] = e
		}
		indexed[name] = m
	}

	sortedKeys := make([]string, 0, len(keys))
	for k := range keys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	rows := make([]Row, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		row := Row{
			Key:    key,
			Values: make(map[string]string, len(envNames)),
		}
		for _, env := range envNames {
			if e, ok := indexed[env][key]; ok {
				row.Values[env] = e.NewVal
				if e.IsDrift() {
					row.Drift = true
				}
			}
		}
		rows = append(rows, row)
	}

	return Table{Envs: envNames, Rows: rows}
}

// DriftOnly returns a filtered Table containing only rows with drift.
func (t Table) DriftOnly() Table {
	filtered := make([]Row, 0)
	for _, r := range t.Rows {
		if r.Drift {
			filtered = append(filtered, r)
		}
	}
	return Table{Envs: t.Envs, Rows: filtered}
}
