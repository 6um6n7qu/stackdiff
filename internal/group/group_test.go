package group_test

import (
	"testing"

	"github.com/yourorg/stackdiff/internal/diff"
	"github.com/yourorg/stackdiff/internal/group"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "DB_HOST", OldVal: "localhost", NewVal: "prod-db", Status: diff.StatusChanged},
		{Key: "DB_PORT", OldVal: "", NewVal: "5432", Status: diff.StatusAdded},
		{Key: "APP_SECRET", OldVal: "abc", NewVal: "", Status: diff.StatusRemoved},
		{Key: "APP_ENV", OldVal: "staging", NewVal: "staging", Status: diff.StatusEqual},
	}
}

func TestApply_GroupsCorrectly(t *testing.T) {
	groups := group.Apply(sampleEntries(), func(e diff.Entry) string {
		return string(e.Status)
	})
	if len(groups) != 4 {
		t.Fatalf("expected 4 groups, got %d", len(groups))
	}
}

func TestApply_SortedByName(t *testing.T) {
	groups := group.ByStatus(sampleEntries())
	for i := 1; i < len(groups); i++ {
		if groups[i].Name < groups[i-1].Name {
			t.Errorf("groups not sorted: %s before %s", groups[i-1].Name, groups[i].Name)
		}
	}
}

func TestByPrefix_GroupsOnUnderscore(t *testing.T) {
	groups := group.ByPrefix(sampleEntries())
	// Expect two prefix groups: DB and APP
	if len(groups) != 2 {
		t.Fatalf("expected 2 prefix groups, got %d", len(groups))
	}
	if groups[0].Name != "APP" {
		t.Errorf("expected first group APP, got %s", groups[0].Name)
	}
}

func TestGroup_HasDrift_True(t *testing.T) {
	g := group.Group{
		Name: "test",
		Entries: []diff.Entry{
			{Key: "X", Status: diff.StatusChanged},
		},
	}
	if !g.HasDrift() {
		t.Error("expected HasDrift true")
	}
}

func TestGroup_HasDrift_False(t *testing.T) {
	g := group.Group{
		Name: "test",
		Entries: []diff.Entry{
			{Key: "X", Status: diff.StatusEqual},
		},
	}
	if g.HasDrift() {
		t.Error("expected HasDrift false")
	}
}

func TestGroup_Count(t *testing.T) {
	g := group.Group{Name: "x", Entries: sampleEntries()}
	if g.Count() != 4 {
		t.Errorf("expected count 4, got %d", g.Count())
	}
}
