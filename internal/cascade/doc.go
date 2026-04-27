// Package cascade provides a multi-stage entry pipeline where each named
// stage transforms the diff entry slice produced by the previous one.
//
// Usage:
//
//	stages := []cascade.Stage{
//		{Name: "filter",    Fn: myFilterFn},
//		{Name: "redact",    Fn: myRedactFn},
//		{Name: "normalize", Fn: myNormalizeFn},
//	}
//	report, err := cascade.Run(entries, stages)
//	if err != nil { ... }
//	fmt.Println(report.HasDrift())
//
// Each stage records how many entries were dropped, which is useful for
// debugging pipeline behaviour without adding global state.
package cascade
