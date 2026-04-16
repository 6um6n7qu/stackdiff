// Package redact masks sensitive configuration values such as passwords,
// tokens, and API keys before they are displayed in diffs, reports, or
// audit logs.
//
// Usage:
//
//	r := redact.New()
//	safe := r.Apply(configMap)
//
// Custom patterns can be added to r.Patterns before calling Apply.
package redact
