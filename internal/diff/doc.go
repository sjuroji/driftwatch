// Package diff computes field-level differences between two consecutive
// drift.Report snapshots produced by the detector.
//
// Typical usage:
//
//	prev, _ := history.Store.List()
//	next := detector.Detect(manifests, liveStates)
//
//	diffs := diff.Compare(prev[0], next)
//	summary := diff.Summarize(diffs)
//	fmt.Println(summary)
//
// The package is intentionally free of I/O and external dependencies so that
// it can be used in tests and in the scheduler loop without side effects.
package diff
