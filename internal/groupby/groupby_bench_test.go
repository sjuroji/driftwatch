package groupby_test

import (
	"fmt"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/groupby"
)

func largeMixedResults(n int) []drift.Result {
	envs := []string{"prod", "staging", "dev", "canary"}
	results := make([]drift.Result, n)
	for i := range results {
		status := drift.StatusInSync
		if i%3 == 0 {
			status = drift.StatusDrifted
		}
		results[i] = drift.Result{
			Service: fmt.Sprintf("svc-%d", i),
			Status:  status,
			Labels:  map[string]string{"env": envs[i%len(envs)]},
		}
	}
	return results
}

func BenchmarkGroup_ByLabel(b *testing.B) {
	results := largeMixedResults(1000)
	fn := groupby.ByLabel("env", "unknown")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = groupby.Group(results, fn)
	}
}

func BenchmarkGroup_ByStatus(b *testing.B) {
	results := largeMixedResults(1000)
	fn := groupby.ByStatus()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = groupby.Group(results, fn)
	}
}

func BenchmarkSummarise(b *testing.B) {
	results := largeMixedResults(1000)
	fn := groupby.ByLabel("env", "unknown")
	groups := groupby.Group(results, fn)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = groupby.Summarise(groups)
	}
}
