package transform_test

import (
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/transform"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "DB_HOST", OldVal: " localhost ", NewVal: " prod.db ", Status: diff.StatusChanged},
		{Key: "API_KEY", OldVal: "", NewVal: "abc123", Status: diff.StatusAdded},
	}
}

func TestLowercaseKeys(t *testing.T) {
	tr := transform.New(transform.LowercaseKeys())
	out := tr.Apply(sampleEntries())
	if out[0].Key != "db_host" {
		t.Errorf("expected db_host, got %s", out[0].Key)
	}
	if out[1].Key != "api_key" {
		t.Errorf("expected api_key, got %s", out[1].Key)
	}
}

func TestTrimSpace(t *testing.T) {
	tr := transform.New(transform.TrimSpace())
	out := tr.Apply(sampleEntries())
	if out[0].OldVal != "localhost" {
		t.Errorf("expected 'localhost', got %q", out[0].OldVal)
	}
	if out[0].NewVal != "prod.db" {
		t.Errorf("expected 'prod.db', got %q", out[0].NewVal)
	}
}

func TestPrefixKey(t *testing.T) {
	tr := transform.New(transform.PrefixKey("svc_"))
	out := tr.Apply(sampleEntries())
	if out[0].Key != "svc_DB_HOST" {
		t.Errorf("unexpected key: %s", out[0].Key)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	orig := sampleEntries()
	tr := transform.New(transform.LowercaseKeys())
	tr.Apply(orig)
	if orig[0].Key != "DB_HOST" {
		t.Error("original entries were mutated")
	}
}

func TestApply_MultipleChained(t *testing.T) {
	tr := transform.New(transform.TrimSpace(), transform.LowercaseKeys())
	out := tr.Apply(sampleEntries())
	if out[0].Key != "db_host" {
		t.Errorf("expected db_host, got %s", out[0].Key)
	}
	if out[0].OldVal != "localhost" {
		t.Errorf("expected localhost, got %q", out[0].OldVal)
	}
}
