// Package score provides drift severity scoring for stackdiff.
//
// It assigns numeric weights to each type of diff entry (changed, added,
// removed) and aggregates them into a Result with a total score and letter
// grade. This allows operators to quickly gauge how significant the detected
// configuration drift is between two environments.
//
// Example:
//
//	res := score.Compute(entries)
//	fmt.Println(res.Grade())  // "A", "B", "C", "D", or "F"
package score
