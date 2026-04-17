package tag_test

import (
	"testing"

	"github.com/you/stackdiff/internal/diff"
	"github.com/you/stackdiff/internal/tag"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "db_host", OldVal: "localhost", NewVal: "prod-db", Status: diff.StatusChanged},
		{Key: "aws_region", OldVal: "", NewVal: "us-east-1", Status: diff.StatusAdded},
		{Key: "log_level", OldVal: "debug", NewVal: "", Status: diff.StatusRemoved},
		{Key: "app_name", OldVal: "svc", NewVal: "svc", Status: diff.StatusEqual},
	}
}

func TestApply_MatchesPrefix(t *testing.T) {
	tagger := tag.New(tag.DefaultRules())
	labels := tagger.Apply(sampleEntries())

	if got := labels["db_host"]; len(got) == 0 || got[0] != "database" {
		t.Errorf("expected db_host to be tagged 'database', got %v", got)
	}
	if got := labels["aws_region"]; len(got) == 0 || got[0] != "cloud" {
		t.Errorf("expected aws_region to be tagged 'cloud', got %v", got)
	}
	if got := labels["log_level"]; len(got) == 0 || got[0] != "logging" {
		t.Errorf("expected log_level to be tagged 'logging', got %v", got)
	}
}

func TestApply_NoMatchReturnsEmpty(t *testing.T) {
	tagger := tag.New(tag.DefaultRules())
	labels := tagger.Apply(sampleEntries())

	if _, ok := labels["app_name"]; ok {
		t.Error("expected app_name to have no tags")
	}
}

func TestApply_CustomRule(t *testing.T) {
	rules := []tag.Rule{{Prefix: "feature_", Tag: "feature-flag"}}
	entries := []diff.Entry{
		{Key: "feature_dark_mode", OldVal: "false", NewVal: "true", Status: diff.StatusChanged},
	}
	tagger := tag.New(rules)
	labels := tagger.Apply(entries)

	if got := labels["feature_dark_mode"]; len(got) == 0 || got[0] != "feature-flag" {
		t.Errorf("expected feature-flag tag, got %v", got)
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	tagger := tag.New(tag.DefaultRules())
	labels := tagger.Apply([]diff.Entry{})
	if len(labels) != 0 {
		t.Errorf("expected empty result, got %v", labels)
	}
}

func TestApply_CaseInsensitivePrefix(t *testing.T) {
	rules := []tag.Rule{{Prefix: "DB_", Tag: "database"}}
	entries := []diff.Entry{
		{Key: "db_password", OldVal: "x", NewVal: "y", Status: diff.StatusChanged},
	}
	tagger := tag.New(rules)
	labels := tagger.Apply(entries)
	if got := labels["db_password"]; len(got) == 0 || got[0] != "database" {
		t.Errorf("expected case-insensitive match, got %v", got)
	}
}
