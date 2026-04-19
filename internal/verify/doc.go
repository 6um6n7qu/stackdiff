// Package verify provides threshold-aware verification of diff entries.
//
// It evaluates a slice of diff.Entry values against a Config and returns
// a Result indicating pass, warning, or fail status along with diagnostic
// messages for each drifted key.
//
// Example:
//
//	cfg := verify.DefaultConfig()
//	cfg.IgnoreAdded = true
//	result := verify.Evaluate(entries, cfg)
//	if !result.Passed() {
//	    fmt.Println(result.Messages)
//	}
package verify
