package rollup_test

import (
	"testing"

	"github.com/yourorg/stackdiff/internal/diff"
	"github.com/yourorg/stackdiff/internal/rollup"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "db.host", OldVal: "localhost", NewVal: "prod-db", Status: diff.StatusChanged},
		{Key: "db.port", OldVal: "", NewVal: "5432", Status: diff.StatusAdded},
		{Key: "cache.host", OldVal: "redis", NewVal: "", Status: diff.StatusRemoved},
		{Key: "cache.ttl", OldVal: "60", NewVal: "60", Status: diff.StatusEqual},
		{Key: "app.debug", OldVal: "true", NewVal: "false", Status: diff.StatusChanged},
	}
}

func TestByKeyFunc_GroupsCorrectly(t *testing.T) {
	groups := rollup.ByKeyFunc(sampleEntries(), rollup.PrefixKeyFunc)
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
}

func TestByKeyFunc_Counts(t *testing.T) {
	groups := rollup.ByKeyFunc(sampleEntries(), rollup.PrefixKeyFunc)
	groupMap := map[string]rollup.Group{}
	for _, g := range groups {
		groupMap[g.Key] = g
	}

	db := groupMap["db"]
	if db.Changed != 1 || db.Added != 1 || db.Removed != 0 {
		t.Errorf("db group counts wrong: %+v", db)
	}

	cache := groupMap["cache"]
	if cache.Removed != 1 || cache.Added != 0 {
		t.Errorf("cache group counts wrong: %+v", cache)
	}
}

func TestByKeyFunc_HasDrift(t *testing.T) {
	groups := rollup.ByKeyFunc(sampleEntries(), rollup.PrefixKeyFunc)
	for _, g := range groups {
		if !g.HasDrift() {
			t.Errorf("expected group %q to have drift", g.Key)
		}
	}
}

func TestByKeyFunc_NoDrift(t *testing.T) {
	entries := []diff.Entry{
		{Key: "app.name", OldVal: "foo", NewVal: "foo", Status: diff.StatusEqual},
	}
	groups := rollup.ByKeyFunc(entries, rollup.PrefixKeyFunc)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group")
	}
	if groups[0].HasDrift() {
		t.Error("expected no drift")
	}
}

func TestByKeyFunc_EmptyEntries(t *testing.T) {
	groups := rollup.ByKeyFunc(nil, rollup.PrefixKeyFunc)
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestByKeyFunc_CustomKeyFunc(t *testing.T) {
	entries := []diff.Entry{
		{Key: "X_DB_HOST", Status: diff.StatusAdded},
		{Key: "X_DB_PORT", Status: diff.StatusAdded},
		{Key: "Y_APP", Status: diff.StatusChanged},
	}
	fn := func(k string) string {
		if len(k) > 0 {
			return string(k[0])
		}
		return "?"
	}
	groups := rollup.ByKeyFunc(entries, fn)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}
