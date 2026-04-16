// Package profile provides named environment profile management for stackdiff.
//
// A Profile is a named configuration that references an env file and carries
// optional labels (e.g. "env=staging", "region=us-east-1"). Profiles are
// stored as YAML files in a directory and can be loaded, listed, and resolved
// into concrete key/value maps ready for drift comparison.
//
// Typical usage:
//
//	store := profile.NewStore("~/.stackdiff/profiles")
//	p, _ := store.Load("staging")
//	rp, _ := profile.Resolve(p)
//	// rp.Entries is map[string]string ready for diff.Compare
package profile
