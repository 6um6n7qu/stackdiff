package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/user/stackdiff/internal/diff"
)

const markdownTmpl = `# StackDiff Report

**Source:** {{ .SourceLabel }}  
**Target:** {{ .TargetLabel }}  
**Generated:** {{ .GeneratedAt.Format "2006-01-02 15:04:05 UTC" }}

## Summary

| Added | Removed | Changed | Total |
|-------|---------|---------|-------|
| {{ .Summary.Added }} | {{ .Summary.Removed }} | {{ .Summary.Changed }} | {{ .Summary.Total }} |

## Changes

{{ range .Diffs }}{{ renderRow . }}
{{ end }}`

// Render writes the report to w in the requested format.
func Render(w io.Writer, r *Report, format Format) error {
	switch format {
	case FormatJSON:
		return renderJSON(w, r)
	case FormatMarkdown:
		return renderMarkdown(w, r)
	default:
		return renderText(w, r)
	}
}

func renderJSON(w io.Writer, r *Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func renderText(w io.Writer, r *Report) error {
	fmt.Fprintf(w, "Source: %s | Target: %s\n", r.SourceLabel, r.TargetLabel)
	fmt.Fprintf(w, "Diffs — added: %d, removed: %d, changed: %d\n\n",
		r.Summary.Added, r.Summary.Removed, r.Summary.Changed)
	for _, d := range r.Diffs {
		fmt.Fprintln(w, formatEntry(d))
	}
	return nil
}

func renderMarkdown(w io.Writer, r *Report) error {
	funcMap := template.FuncMap{
		"renderRow": func(d diff.DiffEntry) string { return "- " + formatEntry(d) },
	}
	tmpl, err := template.New("md").Funcs(funcMap).Parse(markdownTmpl)
	if err != nil {
		return fmt.Errorf("report: markdown template parse: %w", err)
	}
	return tmpl.Execute(w, r)
}

func formatEntry(d diff.DiffEntry) string {
	switch d.Kind {
	case diff.Added:
		return fmt.Sprintf("[+] %s = %q", d.Key, d.NewValue)
	case diff.Removed:
		return fmt.Sprintf("[-] %s = %q", d.Key, d.OldValue)
	case diff.Changed:
		return fmt.Sprintf("[~] %s: %q → %q", d.Key, d.OldValue, d.NewValue)
	default:
		return strings.TrimSpace(fmt.Sprintf("[?] %s", d.Key))
	}
}
