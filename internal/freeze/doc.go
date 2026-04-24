// Package freeze detects violations of frozen configuration keys.
//
// A frozen key is one whose value must remain constant and must not drift
// between environments or across time. Use Enforce to check a set of diff
// entries against a configured list of frozen keys, and MustPass to
// integrate the check into a pipeline that returns an error on any breach.
package freeze
