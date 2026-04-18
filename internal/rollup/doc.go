// Package rollup provides helpers to aggregate drift entries into logical
// groups, making it easier to reason about environment drift at a higher
// level of abstraction (e.g. per-service or per-namespace).
//
// Usage:
//
//	groups := rollup.ByKeyFunc(entries, rollup.PrefixKeyFunc)
//	for _, g := range groups {
//		fmt.Printf("%s: +%d -%d ~%d\n", g.Key, g.Added, g.Removed, g.Changed)
//	}
package rollup
