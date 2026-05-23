// Package notify provides alerting capabilities when config drift is detected.
package notify

import (
	"fmt"
	"io"
	"os"

	"github.com/driftwatch/internal/drift"
)

// Level represents the severity threshold for notifications.
type Level string

const (
	LevelAll    Level = "all"
	LevelDrift  Level = "drift"
	LevelNone   Level = "none"
)

// Notifier sends alerts based on drift reports.
type Notifier interface {
	Notify(report drift.Report) error
}

// Config holds configuration for a notifier.
type Config struct {
	Level   Level
	Writer  io.Writer
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Level:  LevelDrift,
		Writer: os.Stderr,
	}
}

// New returns a Notifier for the given config.
func New(cfg Config) Notifier {
	return &stdNotifier{cfg: cfg}
}

type stdNotifier struct {
	cfg Config
}

func (n *stdNotifier) Notify(report drift.Report) error {
	if n.cfg.Level == LevelNone {
		return nil
	}
	if n.cfg.Level == LevelDrift && !report.HasDrift() {
		return nil
	}
	_, err := fmt.Fprintf(n.cfg.Writer, "[driftwatch] %s\n", report.Summary)
	return err
}
