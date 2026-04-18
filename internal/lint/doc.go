// Package lint provides rule-based linting for stackdiff config entries.
//
// Rules inspect individual diff.Entry values and return findings with a
// severity level (warn or error). The default rule set catches common
// problems such as unfilled placeholders, empty values, and suspiciously
// short secrets.
//
// # Severity Levels
//
// Each finding carries one of two severity levels:
//
//   - [SeverityWarn]: the entry is suspicious but may be intentional.
//   - [SeverityError]: the entry is very likely incorrect and should be fixed.
//
// # Usage
//
//	findings := lint.Apply(entries, lint.DefaultRules())
//	for _, f := range findings {
//		fmt.Println(f)
//	}
//
// Custom rules can be passed alongside or instead of the defaults:
//
//	rules := append(lint.DefaultRules(), myRule)
//	findings := lint.Apply(entries, rules)
package lint
