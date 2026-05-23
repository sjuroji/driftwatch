package metric

import (
	"time"

	"driftwatch/internal/drift"
)

// FromReport builds a RunResult from a completed drift.Report and
// the wall-clock duration the detection cycle took.
func FromReport(report drift.Report, duration time.Duration) RunResult {
	var inSync, drifted int
	for _, entry := range report.Entries {
		if entry.Status == drift.StatusInSync {
			inSync++
		} else {
			drifted++
		}
	}
	return RunResult{
		RunAt:         time.Now(),
		TotalServices: len(report.Entries),
		InSync:        inSync,
		Drifted:       drifted,
		Errors:        0,
		DurationMs:    duration.Milliseconds(),
	}
}
