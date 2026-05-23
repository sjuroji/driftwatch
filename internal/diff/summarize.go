package diff

import (
	"fmt"
	"strings"
)

// Summary holds aggregate statistics for a set of diffs.
type Summary struct {
	Total    int
	ByField  map[string]int
	Services []string
}

// Summarize aggregates a slice of FieldDiff values into a Summary.
func Summarize(diffs []FieldDiff) Summary {
	s := Summary{
		Total:   len(diffs),
		ByField: make(map[string]int),
	}

	seen := make(map[string]struct{})
	for _, d := range diffs {
		s.ByField[d.Field]++
		if _, ok := seen[d.Service]; !ok {
			seen[d.Service] = struct{}{}
			s.Services = append(s.Services, d.Service)
		}
	}
	return s
}

// String returns a compact textual summary.
func (s Summary) String() string {
	if s.Total == 0 {
		return "no changes detected"
	}

	parts := make([]string, 0, len(s.ByField))
	for field, count := range s.ByField {
		parts = append(parts, fmt.Sprintf("%s×%d", field, count))
	}
	return fmt.Sprintf("%d change(s) across %d service(s): %s",
		s.Total, len(s.Services), strings.Join(parts, ", "))
}
