package export_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/stackdiff/internal/diff"
	"github.com/user/stackdiff/internal/export"
	"github.com/user/stackdiff/internal/filter"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "DB_HOST", BaseValue: "localhost", HeadValue: "prod-db", Status: diff.StatusChanged},
		{Key: "LOG_LEVEL", BaseValue: "", HeadValue: "debug", Status: diff.StatusAdded},
		{Key: "OLD_FLAG", BaseValue: "true", HeadValue: "", Status: diff.StatusRemoved},
	}
}

func TestRunPipeline_TextToFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.txt")

	err := export.RunPipeline(export.PipelineInput{
		Entries: sampleEntries(),
		Labels:  map[string]string{"env": "prod"},
		Export:  export.Options{Format: export.FormatText, Dest: dest},
	})
	if err != nil {
		t.Fatalf("RunPipeline error: %v", err)
	}
	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(data), "DB_HOST") {
		t.Error("expected DB_HOST in text output")
	}
}

func TestRunPipeline_JSONToFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "report.json")

	err := export.RunPipeline(export.PipelineInput{
		Entries: sampleEntries(),
		Export:  export.Options{Format: export.FormatJSON, Dest: dest},
	})
	if err != nil {
		t.Fatalf("RunPipeline error: %v", err)
	}
	data, _ := os.ReadFile(dest)
	if !strings.Contains(string(data), "\"key\"") {
		t.Error("expected JSON key field in output")
	}
}

func TestRunPipeline_WithFilter(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "filtered.txt")

	err := export.RunPipeline(export.PipelineInput{
		Entries: sampleEntries(),
		Filter:  filter.Options{OnlyChanged: true},
		Export:  export.Options{Format: export.FormatText, Dest: dest},
	})
	if err != nil {
		t.Fatalf("RunPipeline error: %v", err)
	}
	data, _ := os.ReadFile(dest)
	if strings.Contains(string(data), "LOG_LEVEL") {
		t.Error("filtered output should not contain added entry LOG_LEVEL")
	}
}

func TestRunPipeline_InvalidDest(t *testing.T) {
	err := export.RunPipeline(export.PipelineInput{
		Entries: sampleEntries(),
		Export:  export.Options{Format: export.FormatText, Dest: "/dev/null/bad/path.txt"},
	})
	if err == nil {
		t.Fatal("expected error for invalid destination")
	}
}
