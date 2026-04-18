// Package group provides grouping utilities for diff entries.
//
// Entries can be grouped by any attribute using a GroupFunc, or by
// built-in strategies such as ByStatus or ByPrefix.
//
// Example:
//
//	groups := group.ByStatus(entries)
//	for _, g := range groups {
//		fmt.Println(g.Name, g.Count())
//	}
package group
