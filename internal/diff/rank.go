package diff

import (
	"sort"
)

// RankConfig controls how entries are ranked.
type RankConfig struct {
	// MaxResults limits the number of entries returned. 0 means no limit.
	MaxResults int
	// Weights assigns a numeric weight to each status for ranking.
	// Higher weight = higher rank.
	Weights map[Status]int
}

// DefaultRankConfig returns a RankConfig with sensible defaults.
func DefaultRankConfig() RankConfig {
	return RankConfig{
		MaxResults: 0,
		Weights: map[Status]int{
			StatusChanged: 3,
			StatusAdded:   2,
			StatusRemoved: 2,
			StatusEqual:   0,
		},
	}
}

// RankedEntry pairs an Entry with its computed rank score.
type RankedEntry struct {
	Entry Entry
	Score int
}

// Rank sorts entries by their drift weight, highest first.
// Entries with equal weight preserve their original relative order.
// If cfg.MaxResults > 0, only the top N results are returned.
func Rank(entries []Entry, cfg RankConfig) []RankedEntry {
	weights := cfg.Weights
	if weights == nil {
		weights = DefaultRankConfig().Weights
	}

	ranked := make([]RankedEntry, len(entries))
	for i, e := range entries {
		ranked[i] = RankedEntry{Entry: e, Score: weights[e.Status]}
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})

	if cfg.MaxResults > 0 && len(ranked) > cfg.MaxResults {
		ranked = ranked[:cfg.MaxResults]
	}

	return ranked
}

// RankEntries is a convenience wrapper that returns only the entries
// (without scores) in ranked order.
func RankEntries(entries []Entry, cfg RankConfig) []Entry {
	ranked := Rank(entries, cfg)
	out := make([]Entry, len(ranked))
	for i, r := range ranked {
		out[i] = r.Entry
	}
	return out
}
