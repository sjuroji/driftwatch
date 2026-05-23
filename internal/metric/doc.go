// Package metric collects lightweight run-time statistics for each
// driftwatch detection cycle.
//
// Usage:
//
//	col := metric.New()
//
//	// after each detection run:
//	result := metric.FromReport(report, elapsed)
//	col.Record(result)
//
//	// inspect the latest run:
//	if r, ok := col.Latest(); ok {
//		fmt.Printf("drifted: %d/%d\n", r.Drifted, r.TotalServices)
//	}
//
//	// aggregate across all runs:
//	summary := col.Summarize()
//	fmt.Printf("avg duration: %dms\n", summary.AvgDurationMs)
package metric
