package drift

import (
	"fmt"
	"io"
	"strings"
)

// Report summarises drift results across multiple services.
type Report struct {
	Results []Result
}

// Add appends a drift result to the report.
func (r *Report) Add(result Result) {
	r.Results = append(r.Results, result)
}

// HasDrift returns true if any service in the report has drifted.
func (r *Report) HasDrift() bool {
	for _, res := range r.Results {
		if res.Status == StatusDrifted {
			return true
		}
	}
	return false
}

// Write renders the report as human-readable text to the given writer.
func (r *Report) Write(w io.Writer) error {
	if len(r.Results) == 0 {
		_, err := fmt.Fprintln(w, "No services checked.")
		return err
	}

	for _, res := range r.Results {
		switch res.Status {
		case StatusInSync:
			_, err := fmt.Fprintf(w, "[OK]      %s\n", res.ServiceName)
			if err != nil {
				return err
			}
		case StatusDrifted:
			_, err := fmt.Fprintf(w, "[DRIFTED] %s\n", res.ServiceName)
			if err != nil {
				return err
			}
			for _, diff := range res.Diffs {
				_, err := fmt.Fprintf(w, "          - %s\n", diff)
				if err != nil {
					return err
				}
			}
		default:
			_, err := fmt.Fprintf(w, "[UNKNOWN] %s\n", res.ServiceName)
			if err != nil {
				return err
			}
		}
	}

	summary := buildSummary(r.Results)
	_, err := fmt.Fprintf(w, "\n%s\n", summary)
	return err
}

func buildSummary(results []Result) string {
	total := len(results)
	drifted := 0
	for _, r := range results {
		if r.Status == StatusDrifted {
			drifted++
		}
	}
	parts := []string{
		fmt.Sprintf("total=%d", total),
		fmt.Sprintf("drifted=%d", drifted),
		fmt.Sprintf("in_sync=%d", total-drifted),
	}
	return "Summary: " + strings.Join(parts, " ")
}
