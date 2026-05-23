// Package alert provides threshold-based alerting for drift metrics.
// It evaluates recorded metrics and emits alerts when drift counts
// exceed configured thresholds.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn  Level = "warn"
	LevelCrit  Level = "critical"
	LevelClear Level = "clear"
)

// Config holds thresholds for triggering alerts.
type Config struct {
	// WarnThreshold triggers a warn-level alert when drifted services >= value.
	WarnThreshold int `yaml:"warn_threshold"`
	// CritThreshold triggers a critical alert when drifted services >= value.
	CritThreshold int `yaml:"crit_threshold"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		WarnThreshold: 1,
		CritThreshold: 5,
	}
}

// Alert represents a triggered alert.
type Alert struct {
	Level     Level
	Message   string
	Drifted   int
	Triggered time.Time
}

// Alerter evaluates metrics and writes alerts to a writer.
type Alerter struct {
	cfg Config
	out io.Writer
}

// New creates a new Alerter with the given config.
// If out is nil, os.Stderr is used.
func New(cfg Config, out io.Writer) *Alerter {
	if out == nil {
		out = os.Stderr
	}
	return &Alerter{cfg: cfg, out: out}
}

// Evaluate checks the number of drifted services and returns an Alert.
// Returns nil when drifted count is below the warn threshold.
func (a *Alerter) Evaluate(drifted int) *Alert {
	var lvl Level
	switch {
	case drifted >= a.cfg.CritThreshold:
		lvl = LevelCrit
	case drifted >= a.cfg.WarnThreshold:
		lvl = LevelWarn
	default:
		return nil
	}
	alert := &Alert{
		Level:     lvl,
		Message:   fmt.Sprintf("%d service(s) are drifted", drifted),
		Drifted:   drifted,
		Triggered: time.Now(),
	}
	fmt.Fprintf(a.out, "[ALERT:%s] %s\n", alert.Level, alert.Message)
	return alert
}
