package pivot_test

import (
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/pivot"
)

func entries(pairs ...string) []diff.Entry {
	var out []diff.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, diff.Entry{
			Key:    pairs[i],
			OldVal: "",
			NewVal: pairs[i+1],
			Status: diff.StatusEqual,
		})
	}
	return out
}

func TestBuild_EnvsSorted(t *testing.T) {
	table := pivot.Build(map[string][]diff.Entry{
		"prod": entries("HOST", "prod.example.com"),
		"dev":  entries("HOST", "dev.example.com"),
	})
	if len(table.Envs) != 2 || table.Envs[0] != "dev" {
		t.Fatalf("expected envs sorted, got %v", table.Envs)
	}
}

func TestBuild_RowsContainValues(t *testing.T) {
	table := pivot.Build(map[string][]diff.Entry{
		"prod": entries("PORT", "8080"),
		"dev":  entries("PORT", "3000"),
	})
	if len(table.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(table.Rows))
	}
	row := table.Rows[0]
	if row.Values["prod"] != "8080" || row.Values["dev"] != "3000" {
		t.Errorf("unexpected values: %v", row.Values)
	}
}

func TestBuild_DriftFlagSet(t *testing.T) {
	table := pivot.Build(map[string][]diff.Entry{
		"prod": {
			{Key: "SECRET", OldVal: "old", NewVal: "new", Status: diff.StatusChanged},
		},
		"dev": entries("SECRET", "dev-secret"),
	})
	if !table.Rows[0].Drift {
		t.Error("expected Drift=true for changed entry")
	}
}

func TestBuild_NoDrift(t *testing.T) {
	table := pivot.Build(map[string][]diff.Entry{
		"prod": entries("KEY", "val"),
		"dev":  entries("KEY", "val"),
	})
	if table.Rows[0].Drift {
		t.Error("expected Drift=false")
	}
}

func TestDriftOnly_FiltersRows(t *testing.T) {
	table := pivot.Build(map[string][]diff.Entry{
		"prod": {
			{Key: "A", OldVal: "", NewVal: "1", Status: diff.StatusAdded},
			{Key: "B", OldVal: "x", NewVal: "x", Status: diff.StatusEqual},
		},
	})
	filtered := table.DriftOnly()
	if len(filtered.Rows) != 1 || filtered.Rows[0].Key != "A" {
		t.Errorf("expected only drifted row, got %v", filtered.Rows)
	}
}

func TestBuild_MissingKeyInEnv(t *testing.T) {
	table := pivot.Build(map[string][]diff.Entry{
		"prod": entries("ONLY_PROD", "yes"),
		"dev":  entries("ONLY_DEV", "yes"),
	})
	if len(table.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(table.Rows))
	}
	for _, row := range table.Rows {
		if row.Key == "ONLY_PROD" && row.Values["dev"] != "" {
			t.Error("expected empty value for missing env key")
		}
	}
}
