package digest_test

import (
	"testing"

	"github.com/yourorg/stackdiff/internal/diff"
	"github.com/yourorg/stackdiff/internal/digest"
)

func makeEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "DB_HOST", OldVal: "localhost", NewVal: "prod-db", Status: diff.StatusChanged},
		{Key: "PORT", OldVal: "", NewVal: "5432", Status: diff.StatusAdded},
		{Key: "SECRET", OldVal: "abc", NewVal: "", Status: diff.StatusRemoved},
	}
}

func TestCompute_NonEmpty(t *testing.T) {
	r := digest.Compute(makeEntries())
	if r.Hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if r.EntryCount != 3 {
		t.Fatalf("expected 3 entries, got %d", r.EntryCount)
	}
}

func TestCompute_Empty(t *testing.T) {
	r := digest.Compute(nil)
	if r.Hash == "" {
		t.Fatal("expected non-empty hash for empty input")
	}
	if r.EntryCount != 0 {
		t.Fatalf("expected 0 entries, got %d", r.EntryCount)
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	a := makeEntries()
	b := []diff.Entry{a[2], a[0], a[1]}

	ra := digest.Compute(a)
	rb := digest.Compute(b)

	if !digest.Equal(ra, rb) {
		t.Errorf("expected equal digests regardless of order: %s vs %s", ra.Hash, rb.Hash)
	}
}

func TestCompute_DifferentEntries(t *testing.T) {
	a := makeEntries()
	b := makeEntries()
	b[0].NewVal = "different-host"

	ra := digest.Compute(a)
	rb := digest.Compute(b)

	if digest.Equal(ra, rb) {
		t.Error("expected different digests for different entries")
	}
}

func TestResult_String(t *testing.T) {
	r := digest.Compute(makeEntries())
	s := r.String()
	if len(s) == 0 {
		t.Error("expected non-empty string representation")
	}
}

func TestEqual_SameInput(t *testing.T) {
	r1 := digest.Compute(makeEntries())
	r2 := digest.Compute(makeEntries())
	if !digest.Equal(r1, r2) {
		t.Error("expected equal digests for identical input")
	}
}
