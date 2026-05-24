package threshold

import "github.com/driftwatch/internal/drift"

// Result pairs a service name with its computed severity Level.
type Result struct {
	Service string
	Level   Level
	Drifted int
}

// ClassifyReport applies the Evaluator to every entry in a drift report
// and returns one Result per entry.
//
// An entry is considered drifted when its Status is not "in-sync".
// The drifted-field count is derived from the number of populated
// Diff fields on the entry (Replicas, Image, Env).
func (e *Evaluator) ClassifyReport(report drift.Report) []Result {
	results := make([]Result, 0, len(report.Entries))
	for _, entry := range report.Entries {
		count := countDriftedFields(entry)
		results = append(results, Result{
			Service: entry.Service,
			Level:   e.Classify(count),
			Drifted: count,
		})
	}
	return results
}

// countDriftedFields returns the number of fields that differ between
// the declared manifest and the live state for a single report entry.
func countDriftedFields(entry drift.Entry) int {
	if entry.Status == "in-sync" {
		return 0
	}
	count := 0
	if entry.ReplicaDrift {
		count++
	}
	if entry.ImageDrift {
		count++
	}
	if entry.EnvDrift {
		count++
	}
	return count
}
