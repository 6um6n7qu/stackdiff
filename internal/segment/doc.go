// Package segment provides utilities for splitting diff entry slices into
// named buckets based on user-defined predicate rules.
//
// Typical usage:
//
//	result, err := segment.Apply(entries, segment.DefaultRules())
//	if err != nil {
//		log.Fatal(err)
//	}
//	for name, bucket := range result.Buckets {
//		fmt.Printf("%s: %d entries\n", name, len(bucket))
//	}
package segment
