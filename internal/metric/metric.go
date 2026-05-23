// Package metric provides lightweight run-time counters and gauges
// that driftwatch records after every detection cycle.
package metric

import (
	"sync"
	"time"
)

// RunResult holds the counters produced by a single detection run.
type RunResult struct {
	RunAt         time.Time
	TotalServices int
	InSync        int
	Drifted       int
	Errors        int
	DurationMs    int64
}

// Collector accumulates RunResults in memory.
type Collector struct {
	mu      sync.Mutex
	results []RunResult
}

// New returns an initialised Collector.
func New() *Collector {
	return &Collector{}
}

// Record appends a RunResult to the collector.
func (c *Collector) Record(r RunResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.results = append(c.results, r)
}

// Latest returns the most recently recorded RunResult and true,
// or a zero value and false when no results have been recorded yet.
func (c *Collector) Latest() (RunResult, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.results) == 0 {
		return RunResult{}, false
	}
	return c.results[len(c.results)-1], true
}

// All returns a copy of every recorded RunResult, oldest first.
func (c *Collector) All() []RunResult {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]RunResult, len(c.results))
	copy(out, c.results)
	return out
}

// Summary aggregates totals across all recorded runs.
type Summary struct {
	Runs        int
	TotalDrifts int
	TotalErrors int
	AvgDurationMs int64
}

// Summarize computes aggregate statistics over all recorded runs.
func (c *Collector) Summarize() Summary {
	c.mu.Lock()
	defer c.mu.Unlock()
	s := Summary{Runs: len(c.results)}
	if s.Runs == 0 {
		return s
	}
	var totalDur int64
	for _, r := range c.results {
		s.TotalDrifts += r.Drifted
		s.TotalErrors += r.Errors
		totalDur += r.DurationMs
	}
	s.AvgDurationMs = totalDur / int64(s.Runs)
	return s
}
