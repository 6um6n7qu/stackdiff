// Package suppress silences known or expected drift entries using configurable
// rules. Rules may target exact keys or key prefixes and can carry an optional
// expiry time so that temporary suppressions are automatically lifted.
//
// Typical usage:
//
//	rules := []suppress.Rule{
//		{Key: "DEPLOY_SHA", Reason: "changes every deploy"},
//		{Key: "TMP_*",     Reason: "temporary vars are noisy"},
//	}
//	s := suppress.New(rules)
//	filtered := s.Apply(entries)
package suppress
