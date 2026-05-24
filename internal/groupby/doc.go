// Package groupby provides utilities for partitioning drift.Result slices
// into named groups using pluggable key functions.
//
// Common use-cases include grouping results by environment label, team label,
// or drift status so that callers can produce per-group summaries or targeted
// notifications without duplicating iteration logic.
//
// Usage:
//
//	fn := groupby.ByLabel("env", "unknown")
//	groups := groupby.Group(report.Results, fn)
//	summaries := groupby.Summarise(groups)
//	for _, s := range summaries {
//		fmt.Printf("%s: %d/%d drifted\n", s.Key, s.Drifted, s.Total)
//	}
package groupby
