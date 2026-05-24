// Package groupby provides utilities for grouping drift reports
// by arbitrary string keys such as environment, team, or region.
package groupby

import (
	"fmt"

	"github.com/yourorg/driftwatch/internal/drift"
)

// KeyFunc extracts a grouping key from a drift result.
type KeyFunc func(r drift.Result) string

// ByLabel returns a KeyFunc that groups results by a specific label key.
// Results missing the label are placed under the provided fallback string.
func ByLabel(labelKey, fallback string) KeyFunc {
	return func(r drift.Result) string {
		if v, ok := r.Labels[labelKey]; ok && v != "" {
			return v
		}
		return fallback
	}
}

// ByStatus returns a KeyFunc that groups results by their drift status string.
func ByStatus() KeyFunc {
	return func(r drift.Result) string {
		return r.Status
	}
}

// Group partitions a slice of drift.Result into a map keyed by the value
// returned from fn. The original order within each group is preserved.
func Group(results []drift.Result, fn KeyFunc) map[string][]drift.Result {
	out := make(map[string][]drift.Result)
	for _, r := range results {
		k := fn(r)
		out[k] = append(out[k], r)
	}
	return out
}

// Summary holds aggregate counts for a single group.
type Summary struct {
	Key      string
	Total    int
	Drifted  int
	InSync   int
}

// Summarise converts a grouped map into a slice of Summary values.
func Summarise(groups map[string][]drift.Result) []Summary {
	out := make([]Summary, 0, len(groups))
	for key, results := range groups {
		s := Summary{Key: key, Total: len(results)}
		for _, r := range results {
			if r.Status == drift.StatusDrifted {
				s.Drifted++
			} else {
				s.InSync++
			}
		}
		out = append(out, s)
	}
	return out
}

// Validate ensures fn is not nil.
func Validate(fn KeyFunc) error {
	if fn == nil {
		return fmt.Errorf("groupby: KeyFunc must not be nil")
	}
	return nil
}
