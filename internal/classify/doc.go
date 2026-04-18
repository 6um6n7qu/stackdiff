// Package classify assigns severity levels (critical, high, medium, low) to drift
// entries based on configurable rules. Rules are evaluated in order and the first
// match wins, making rule ordering significant.
//
// Example usage:
//
//	c := classify.New(classify.DefaultRules())
//	results := c.Apply(entries)
//	for _, r := range results {
//		fmt.Println(r.Entry.Key, r.Severity)
//	}
package classify
