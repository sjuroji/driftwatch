// Package rollup aggregates drift reports across multiple runs into
// a summarised view, making it easy to spot persistently drifted services.
package rollup

import (
	"fmt"
	"sort"

	"github.com/driftwatch/internal/drift"
)

// ServiceSummary holds aggregated drift statistics for a single service
// across all reports that were rolled up.
type ServiceSummary struct {
	Service    string
	TotalRuns  int
	DriftCount int
	// DriftRate is the fraction of runs in which this service was drifted.
	DriftRate float64
	LastStatus string
}

// Result is the output of a rollup operation.
type Result struct {
	TotalReports int
	Services     []ServiceSummary
}

// Aggregate accepts a slice of drift reports (oldest first) and returns a
// Result that summarises per-service drift frequency.
func Aggregate(reports []drift.Report) (Result, error) {
	if len(reports) == 0 {
		return Result{}, fmt.Errorf("rollup: no reports provided")
	}

	type counters struct {
		total      int
		driftCount int
		lastStatus string
	}

	tracker := make(map[string]*counters)

	for _, report := range reports {
		for _, entry := range report.Entries {
			c, ok := tracker[entry.Service]
			if !ok {
				c = &counters{}
				tracker[entry.Service] = c
			}
			c.total++
			if entry.Status == drift.StatusDrift {
				c.driftCount++
			}
			c.lastStatus = string(entry.Status)
		}
	}

	summaries := make([]ServiceSummary, 0, len(tracker))
	for svc, c := range tracker {
		rate := 0.0
		if c.total > 0 {
			rate = float64(c.driftCount) / float64(c.total)
		}
		summaries = append(summaries, ServiceSummary{
			Service:    svc,
			TotalRuns:  c.total,
			DriftCount: c.driftCount,
			DriftRate:  rate,
			LastStatus: c.lastStatus,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		if summaries[i].DriftRate != summaries[j].DriftRate {
			return summaries[i].DriftRate > summaries[j].DriftRate
		}
		return summaries[i].Service < summaries[j].Service
	})

	return Result{
		TotalReports: len(reports),
		Services:     summaries,
	}, nil
}
