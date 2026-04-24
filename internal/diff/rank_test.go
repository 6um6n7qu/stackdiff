package diff

import (
	"testing"
)

func makeRankEntries() []Entry {
	return []Entry{
		{Key: "A", OldValue: "1", NewValue: "1", Status: StatusEqual},
		{Key: "B", OldValue: "x", NewValue: "y", Status: StatusChanged},
		{Key: "C", OldValue: "", NewValue: "z", Status: StatusAdded},
		{Key: "D", OldValue: "w", NewValue: "", Status: StatusRemoved},
	}
}

func TestRank_ChangedFirst(t *testing.T) {
	entries := makeRankEntries()
	ranked := Rank(entries, DefaultRankConfig())

	if ranked[0].Entry.Key != "B" {
		t.Errorf("expected Changed entry first, got %s", ranked[0].Entry.Key)
	}
}

func TestRank_EqualLast(t *testing.T) {
	entries := makeRankEntries()
	ranked := Rank(entries, DefaultRankConfig())

	last := ranked[len(ranked)-1]
	if last.Entry.Status != StatusEqual {
		t.Errorf("expected Equal entry last, got status %v", last.Entry.Status)
	}
}

func TestRank_ScoreAssigned(t *testing.T) {
	entries := makeRankEntries()
	ranked := Rank(entries, DefaultRankConfig())

	for _, r := range ranked {
		switch r.Entry.Status {
		case StatusChanged:
			if r.Score != 3 {
				t.Errorf("Changed: expected score 3, got %d", r.Score)
			}
		case StatusAdded, StatusRemoved:
			if r.Score != 2 {
				t.Errorf("%v: expected score 2, got %d", r.Entry.Status, r.Score)
			}
		case StatusEqual:
			if r.Score != 0 {
				t.Errorf("Equal: expected score 0, got %d", r.Score)
			}
		}
	}
}

func TestRank_MaxResults(t *testing.T) {
	entries := makeRankEntries()
	cfg := DefaultRankConfig()
	cfg.MaxResults = 2

	ranked := Rank(entries, cfg)
	if len(ranked) != 2 {
		t.Errorf("expected 2 results, got %d", len(ranked))
	}
}

func TestRank_MaxResultsZeroReturnsAll(t *testing.T) {
	entries := makeRankEntries()
	cfg := DefaultRankConfig()
	cfg.MaxResults = 0

	ranked := Rank(entries, cfg)
	if len(ranked) != len(entries) {
		t.Errorf("expected %d results, got %d", len(entries), len(ranked))
	}
}

func TestRankEntries_ReturnsOnlyEntries(t *testing.T) {
	entries := makeRankEntries()
	out := RankEntries(entries, DefaultRankConfig())

	if len(out) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(out))
	}
	// First should still be the Changed entry
	if out[0].Status != StatusChanged {
		t.Errorf("expected first entry to be Changed, got %v", out[0].Status)
	}
}

func TestRank_StableOrder_EqualWeights(t *testing.T) {
	entries := []Entry{
		{Key: "X", OldValue: "", NewValue: "1", Status: StatusAdded},
		{Key: "Y", OldValue: "2", NewValue: "", Status: StatusRemoved},
	}
	ranked := Rank(entries, DefaultRankConfig())

	// Both have weight 2; original order should be preserved (stable sort)
	if ranked[0].Entry.Key != "X" || ranked[1].Entry.Key != "Y" {
		t.Errorf("expected stable order X,Y; got %s,%s", ranked[0].Entry.Key, ranked[1].Entry.Key)
	}
}
