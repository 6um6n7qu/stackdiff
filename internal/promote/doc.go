// Package promote implements environment promotion for stackdiff.
//
// A promotion compares a source snapshot (e.g. staging) against a target
// snapshot (e.g. production) and produces the patch operations needed to
// bring the target in line with the source.
//
// Basic usage:
//
//	res, err := promote.Run(stagingSnap, prodSnap)
//	if err != nil { ... }
//	fmt.Println(res.Summary())
//
// Use RunWithConfig together with functional options for dry-run support
// and key exclusion:
//
//	cfg := promote.NewConfig(
//		promote.WithDryRun(true),
//		promote.WithIgnoreKeys("DB_PASSWORD", "API_SECRET"),
//	)
//	res, err := promote.RunWithConfig(stagingSnap, prodSnap, cfg)
package promote
