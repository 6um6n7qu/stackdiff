// Package label attaches structured key=value metadata to diff entries
// using a rule-based system.
//
// Each Rule defines a match predicate and a label producer. When a rule
// matches an entry, the produced labels are merged into the entry's Meta
// map. Multiple rules can contribute labels to the same entry.
//
// Example usage:
//
//	l := label.New(label.DefaultRules())
//	annotated := l.Apply(entries)
//	for _, e := range annotated {
//		fmt.Println(e.Key, label.Format(e.Meta))
//	}
//
// DefaultRules provides status and sensitive-key labeling out of the box.
// Custom rules can be composed alongside or instead of the defaults.
package label
