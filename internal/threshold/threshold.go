// Package threshold provides configurable drift severity thresholds
// that classify a drift report entry as low, medium, or high severity
// based on the number and type of drifted fields.
package threshold

import "fmt"

// Level represents a drift severity level.
type Level string

const (
	LevelNone   Level = "none"
	LevelLow    Level = "low"
	LevelMedium Level = "medium"
	LevelHigh   Level = "high"
)

// Config holds the field-count boundaries for each severity level.
type Config struct {
	// LowAt is the minimum number of drifted fields to classify as low.
	LowAt int `yaml:"low_at"`
	// MediumAt is the minimum number of drifted fields to classify as medium.
	MediumAt int `yaml:"medium_at"`
	// HighAt is the minimum number of drifted fields to classify as high.
	HighAt int `yaml:"high_at"`
}

// DefaultConfig returns sensible out-of-the-box thresholds.
func DefaultConfig() Config {
	return Config{
		LowAt:    1,
		MediumAt: 3,
		HighAt:   6,
	}
}

// Evaluator classifies drift counts into severity levels.
type Evaluator struct {
	cfg Config
}

// New creates an Evaluator from the supplied Config.
// Returns an error if the boundaries are not strictly ascending.
func New(cfg Config) (*Evaluator, error) {
	if cfg.LowAt <= 0 || cfg.MediumAt <= cfg.LowAt || cfg.HighAt <= cfg.MediumAt {
		return nil, fmt.Errorf(
			"threshold: boundaries must satisfy 0 < low(%d) < medium(%d) < high(%d)",
			cfg.LowAt, cfg.MediumAt, cfg.HighAt,
		)
	}
	return &Evaluator{cfg: cfg}, nil
}

// Classify returns the Level that corresponds to driftedFields.
func (e *Evaluator) Classify(driftedFields int) Level {
	switch {
	case driftedFields <= 0:
		return LevelNone
	case driftedFields < e.cfg.MediumAt:
		return LevelLow
	case driftedFields < e.cfg.HighAt:
		return LevelMedium
	default:
		return LevelHigh
	}
}
