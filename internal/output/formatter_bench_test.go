package output_test

import (
	"io"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/output"
)

func BenchmarkTextFormatter(b *testing.B) {
	f, _ := output.New(output.FormatText)
	report := sampleReport()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Write(io.Discard, report) //nolint:errcheck
	}
}

func BenchmarkJSONFormatter(b *testing.B) {
	f, _ := output.New(output.FormatJSON)
	report := sampleReport()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Write(io.Discard, report) //nolint:errcheck
	}
}

func BenchmarkTextFormatter_LargeReport(b *testing.B) {
	f, _ := output.New(output.FormatText)
	report := drift.Report{
		Summary: drift.Summary{TotalServices: 50, InSync: 25, Drifted: 25},
	}
	for i := 0; i < 50; i++ {
		report.Entries = append(report.Entries, drift.Entry{
			ServiceName: "service",
			InSync:      i%2 == 0,
			Diffs: []drift.Diff{
				{Field: "replicas", Want: 3, Got: 1},
				{Field: "image", Want: "v1.0", Got: "v0.9"},
			},
		})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Write(io.Discard, report) //nolint:errcheck
	}
}
