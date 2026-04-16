package history

import (
	"sort"
	"time"
)

// QueryOptions holds filtering criteria for history records.
type QueryOptions struct {
	Since    *time.Time
	Until    *time.Time
	OnlyDrift bool
	Limit    int
}

// Query filters records in the store according to the provided options.
// Results are returned in reverse chronological order (newest first).
func (s *Store) Query(opts QueryOptions) []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []Record
	for _, r := range s.records {
		if opts.OnlyDrift && !r.HasDrift() {
			continue
		}
		if opts.Since != nil && r.Timestamp.Before(*opts.Since) {
			continue
		}
		if opts.Until != nil && r.Timestamp.After(*opts.Until) {
			continue
		}
		results = append(results, r)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.After(results[j].Timestamp)
	})

	if opts.Limit > 0 && len(results) > opts.Limit {
		results = results[:opts.Limit]
	}

	return results
}

// Latest returns the most recent record, or false if the store is empty.
func (s *Store) Latest() (Record, bool) {
	results := s.Query(QueryOptions{Limit: 1})
	if len(results) == 0 {
		return Record{}, false
	}
	return results[0], true
}
