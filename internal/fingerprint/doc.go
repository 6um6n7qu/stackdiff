// Package fingerprint provides stable SHA-256 based identity hashing for
// config maps and diff entry slices.
//
// Use Compute to hash a raw key-value map, or FromEntries to hash a slice
// of diff.Entry values. Equal compares two Result values for identity.
package fingerprint
