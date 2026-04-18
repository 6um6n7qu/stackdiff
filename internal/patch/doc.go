// Package patch provides utilities for generating and applying minimal
// configuration patches between two key-value maps.
//
// A patch is an ordered list of Ops (set or delete) that, when applied
// to the source map, produces the destination map. This is useful for
// generating actionable remediation steps after a drift comparison.
package patch
