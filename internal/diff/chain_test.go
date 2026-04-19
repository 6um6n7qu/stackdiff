package diff

import (
	"strings"
	"testing"
)

func makeEntries() []Entry {
	return []Entry{
		{Key: "HOST", OldVal: "localhost", NewVal: "prod.example.com", Status: StatusChanged},
		{Key: "PORT", OldVal: "8080", NewVal: "8080", Status: StatusEqual},
		{Key: "SECRET", OldVal: "", NewVal: "abc123", Status: StatusAdded},
	}
}

func filterChanged(entries []Entry) []Entry {
	var out []Entry
	for _, e := range entries {
		if e.Status != StatusEqual {
			out = append(out, e)
		}
	}
	return out
}

func uppercaseKeys(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		e.Key = strings.ToUpper(e.Key)
		out[i] = e
	}
	return out
}

func TestChain_NoSteps(t *testing.T) {
	entries := makeEntries()
	result := Chain(entries, nil)
	if len(result) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(result))
	}
}

func TestChain_SingleStep(t *testing.T) {
	entries := makeEntries()
	steps := []Step{{Name: "filter", Apply: filterChanged}}
	result := Chain(entries, steps)
	if len(result) != 2 {
		t.Fatalf("expected 2 drift entries, got %d", len(result))
	}
}

func TestChain_MultipleSteps(t *testing.T) {
	entries := makeEntries()
	steps := []Step{
		{Name: "filter", Apply: filterChanged},
		{Name: "uppercase", Apply: uppercaseKeys},
	}
	result := Chain(entries, steps)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key != strings.ToUpper(e.Key) {
			t.Errorf("expected uppercase key, got %q", e.Key)
		}
	}
}

func TestChainWithTrace_RecordsIntermediates(t *testing.T) {
	entries := makeEntries()
	steps := []Step{
		{Name: "filter", Apply: filterChanged},
		{Name: "uppercase", Apply: uppercaseKeys},
	}
	trace := ChainWithTrace(entries, steps)
	if _, ok := trace["filter"]; !ok {
		t.Error("expected 'filter' trace entry")
	}
	if _, ok := trace["uppercase"]; !ok {
		t.Error("expected 'uppercase' trace entry")
	}
	if len(trace["filter"]) != 2 {
		t.Errorf("expected 2 after filter, got %d", len(trace["filter"]))
	}
}

func TestChainWithTrace_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries()
	orig := len(entries)
	steps := []Step{{Name: "filter", Apply: filterChanged}}
	ChainWithTrace(entries, steps)
	if len(entries) != orig {
		t.Error("original entries slice was mutated")
	}
}
