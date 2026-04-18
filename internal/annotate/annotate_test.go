package annotate_test

import (
	"testing"

	"github.com/stackdiff/stackdiff/internal/annotate"
	"github.com/stackdiff/stackdiff/internal/diff"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "APP_PORT", OldVal: "", NewVal: "8080", Status: diff.StatusAdded},
		{Key: "DB_HOST", OldVal: "old", NewVal: "", Status: diff.StatusRemoved},
		{Key: "LOG_LEVEL", OldVal: "info", NewVal: "debug", Status: diff.StatusChanged},
		{Key: "API_TOKEN", OldVal: "x", NewVal: "y", Status: diff.StatusChanged},
		{Key: "REGION", OldVal: "us-east", NewVal: "us-east", Status: diff.StatusEqual},
	}
}

func TestApply_AnnotatesAdded(t *testing.T) {
	a := annotate.New(annotate.DefaultRules())
	anns := a.Apply(sampleEntries())
	for _, ann := range anns {
		if ann.Key == "APP_PORT" {
			if ann.Meta["status"] != "added" {
				t.Errorf("expected added, got %s", ann.Meta["status"])
			}
			return
		}
	}
	t.Error("APP_PORT annotation not found")
}

func TestApply_AnnotatesRemoved(t *testing.T) {
	a := annotate.New(annotate.DefaultRules())
	anns := a.Apply(sampleEntries())
	for _, ann := range anns {
		if ann.Key == "DB_HOST" && ann.Meta["status"] == "removed" {
			return
		}
	}
	t.Error("DB_HOST removed annotation not found")
}

func TestApply_SensitiveKey(t *testing.T) {
	a := annotate.New(annotate.DefaultRules())
	anns := a.Apply(sampleEntries())
	for _, ann := range anns {
		if ann.Key == "API_TOKEN" && ann.Meta["category"] == "sensitive" {
			return
		}
	}
	t.Error("API_TOKEN sensitive annotation not found")
}

func TestApply_EqualNotAnnotated(t *testing.T) {
	a := annotate.New(annotate.DefaultRules())
	anns := a.Apply(sampleEntries())
	for _, ann := range anns {
		if ann.Key == "REGION" {
			t.Error("REGION should not be annotated")
		}
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	a := annotate.New(annotate.DefaultRules())
	anns := a.Apply([]diff.Entry{})
	if len(anns) != 0 {
		t.Errorf("expected 0 annotations, got %d", len(anns))
	}
}

func TestApply_CustomRule(t *testing.T) {
	rules := []annotate.Rule{
		{
			Match:  func(e diff.Entry) bool { return e.Key == "LOG_LEVEL" },
			Labels: map[string]string{"category": "logging"},
		},
	}
	a := annotate.New(rules)
	anns := a.Apply(sampleEntries())
	if len(anns) != 1 || anns[0].Key != "LOG_LEVEL" {
		t.Errorf("expected 1 annotation for LOG_LEVEL, got %d", len(anns))
	}
}
