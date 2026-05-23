package metric_test

import (
	"testing"
	"time"

	"driftwatch/internal/metric"
)

func sampleResult(drifted, errors int, durMs int64) metric.RunResult {
	return metric.RunResult{
		RunAt:         time.Now(),
		TotalServices: 5,
		InSync:        5 - drifted,
		Drifted:       drifted,
		Errors:        errors,
		DurationMs:    durMs,
	}
}

func TestNew_StartsEmpty(t *testing.T) {
	c := metric.New()
	if got := c.All(); len(got) != 0 {
		t.Fatalf("expected empty collector, got %d results", len(got))
	}
}

func TestLatest_NoResults(t *testing.T) {
	c := metric.New()
	_, ok := c.Latest()
	if ok {
		t.Fatal("expected ok=false on empty collector")
	}
}

func TestRecord_And_Latest(t *testing.T) {
	c := metric.New()
	r1 := sampleResult(0, 0, 10)
	r2 := sampleResult(2, 1, 20)
	c.Record(r1)
	c.Record(r2)

	got, ok := c.Latest()
	if !ok {
		t.Fatal("expected ok=true")
	}
	if got.Drifted != 2 {
		t.Errorf("expected Drifted=2, got %d", got.Drifted)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	c := metric.New()
	c.Record(sampleResult(1, 0, 5))
	c.Record(sampleResult(3, 0, 8))

	all := c.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 results, got %d", len(all))
	}
	// mutating the returned slice must not affect the collector
	all[0].Drifted = 99
	if got, _ := c.Latest(); got.Drifted == 99 {
		t.Error("All() returned a reference, not a copy")
	}
}

func TestSummarize_Empty(t *testing.T) {
	c := metric.New()
	s := c.Summarize()
	if s.Runs != 0 || s.TotalDrifts != 0 {
		t.Errorf("unexpected summary on empty collector: %+v", s)
	}
}

func TestSummarize_Aggregates(t *testing.T) {
	c := metric.New()
	c.Record(sampleResult(2, 1, 10))
	c.Record(sampleResult(4, 0, 30))

	s := c.Summarize()
	if s.Runs != 2 {
		t.Errorf("expected Runs=2, got %d", s.Runs)
	}
	if s.TotalDrifts != 6 {
		t.Errorf("expected TotalDrifts=6, got %d", s.TotalDrifts)
	}
	if s.TotalErrors != 1 {
		t.Errorf("expected TotalErrors=1, got %d", s.TotalErrors)
	}
	if s.AvgDurationMs != 20 {
		t.Errorf("expected AvgDurationMs=20, got %d", s.AvgDurationMs)
	}
}
