// Package export provides utilities for exporting diff reports to various
// formats and destinations.
//
// The primary entry point is [RunPipeline], which accepts raw diff entries,
// applies optional filters, builds a [report.Report], renders it in the
// requested format (text, JSON, or Markdown), and writes the result to either
// stdout or a file on disk.
//
// Supported formats:
//
//	"text"     — human-readable columnar output (default)
//	"json"     — machine-readable JSON
//	"markdown" — GitHub-flavoured Markdown table
//
// Destination is controlled via Options.Dest:
//
//	"-"        — write to stdout
//	"<path>"   — write to the given file (parent dirs are created automatically)
package export
