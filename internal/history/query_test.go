package history

import (
	"testing"
	"time"

	"github.com/user/stackdiff/internal/diff"
)

func buildStore(t *testing.T, records []Record) *Store {
	t.Helper()
	s := &Store{}
	s.records = records
	return s
}

func makeRecord(ts time.Time, entries []diff.Entry) Record {
	r := NewRecord("env-a", "env-b", entries)
	r.Timestamp = ts
	return r
}

// sampleEntries returns a small set of diff entries for use in tests
// that require at least one drift entry to be present.
func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "service.replicas", Left: "2", Right: "3"},
	}
}

func TestQuery_NoFilters(t *testing.T) {
	now := time.Now()
	records := []Record{
		makeRecord(now.Add(-2*time.Hour), sampleEntries()),
		makeRecord(now.Add(-1*time.Hour), nil),
	}
	s := buildStore(t, records)
	got := s.Query(QueryOptions{})
	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}
}

func TestQuery_OnlyDrift(t *testing.T) {
	now := time.Now()
	records := []Record{
		makeRecord(now.Add(-2*time.Hour), sampleEntries()),
		makeRecord(now.Add(-1*time.Hour), nil),
	}
	s := buildStore(t, records)
	got := s.Query(QueryOptions{OnlyDrift: true})
	if len(got) != 1 {
		t.Fatalf("expected 1 drift record, got %d", len(got))
	}
}

func TestQuery_Since(t *testing.T) {
	now := time.Now()
	cutoff := now.Add(-90 * time.Minute)
	records := []Record{
		makeRecord(now.Add(-2*time.Hour), nil),
		makeRecord(now.Add(-1*time.Hour), nil),
	}
	s := buildStore(t, records)
	got := s.Query(QueryOptions{Since: &cutoff})
	if len(got) != 1 {
		t.Fatalf("expected 1 record after cutoff, got %d", len(got))
	}
}

func TestQuery_Limit(t *testing.T) {
	now := time.Now()
	records := []Record{
		makeRecord(now.Add(-3*time.Hour), nil),
		makeRecord(now.Add(-2*time.Hour), nil),
		makeRecord(now.Add(-1*time.Hour), nil),
	}
	s := buildStore(t, records)
	got := s.Query(QueryOptions{Limit: 2})
	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}
}

func TestQuery_ReverseChronological(t *testing.T) {
	now := time.Now()
	old := makeRecord(now.Add(-2*time.Hour), nil)
	recent := makeRecord(now.Add(-1*time.Hour), nil)
	s := buildStore(t, []Record{old, recent})
	got := s.Query(QueryOptions{})
	if !got[0].Timestamp.Equal(recent.Timestamp) {
		t.Error("expected newest record first")
	}
}

func TestLatest_ReturnsNewest(t *testing.T) {
	now := time.Now()
	old := makeRecord(now.Add(-2*time.Hour), nil)
	recent := makeRecord(now.Add(-1*time.Hour), sampleEntries())
	s := buildStore(t, []Record{old, recent})
	got, ok := s.Latest()
	if !ok {
		t.Fatal("expected a record")
	}
	if !got.Timestamp.Equal(recent.Timestamp) {
		t.Error("expected latest record")
	}
}

func TestLatest_EmptyStore(t *testing.T) {
	s := buildStore(t, nil)
	_, ok := s.Latest()
	if ok {
		t.Error("expected no record from empty store")
	}
}
