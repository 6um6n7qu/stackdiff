// Package schedule provides a periodic execution engine for stackdiff.
//
// It runs diff comparisons at a configurable interval and invokes
// callbacks when drift is detected or an error occurs.
//
// Example usage:
//
//	schedule.Run(ctx, schedule.Job{
//		Interval: 30 * time.Second,
//		Loader:   myLoader,
//		OnDrift:  handleDrift,
//		OnError:  log.Println,
//	})
package schedule
