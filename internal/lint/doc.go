// Package lint provides rule-based linting for stackdiff config entries.
//
// Rules inspect individual diff.Entry values and return findings with a
// severity level (warn or error). The default rule set catches common
// problems such as unfilled placeholders, empty values, and suspiciously
// short secrets.
//
// Usage:
//
//	findings := lint.Apply(entries, lint.DefaultRules())
//	for _, f := range findings {
//		fmt.Println(f)
//	}
package lint
