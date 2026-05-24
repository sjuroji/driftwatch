// Package severity assigns a severity level to a drift report based on the
// number and type of drifted fields across all services.
package severity

import "github.com/driftwatch/internal/drift"

// Level represents a drift severity classification.
type Level string

const (
	LevelNone     Level = "none"
	LevelLow      Level = "low"
	LevelMedium   Level = "medium"
	LevelHigh     Level = "high"
	LevelCritical Level = "critical"
)

// Config holds the thresholds used to determine severity.
// Thresholds are the minimum number of drifted services required to reach
// the corresponding level.
type Config struct {
	LowAt      int // default 1
	MediumAt   int // default 3
	HighAt     int // default 6
	CriticalAt int // default 10
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		LowAt:      1,
		MediumAt:   3,
		HighAt:     6,
		CriticalAt: 10,
	}
}

// Evaluator assigns severity levels to drift reports.
type Evaluator struct {
	cfg Config
}

// New creates an Evaluator with the provided config.
// If cfg is zero-value, DefaultConfig is used.
func New(cfg Config) *Evaluator {
	if cfg.LowAt == 0 && cfg.MediumAt == 0 && cfg.HighAt == 0 && cfg.CriticalAt == 0 {
		cfg = DefaultConfig()
	}
	return &Evaluator{cfg: cfg}
}

// Evaluate inspects the report and returns the appropriate Level.
func (e *Evaluator) Evaluate(report drift.Report) Level {
	drifted := 0
	for _, entry := range report.Entries {
		if entry.Status == drift.StatusDrifted {
			drifted++
		}
	}

	switch {
	case drifted >= e.cfg.CriticalAt:
		return LevelCritical
	case drifted >= e.cfg.HighAt:
		return LevelHigh
	case drifted >= e.cfg.MediumAt:
		return LevelMedium
	case drifted >= e.cfg.LowAt:
		return LevelLow
	default:
		return LevelNone
	}
}
