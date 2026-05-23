// Package audit provides structured audit logging for drift detection events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an audit event.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	Message   string            `json:"message"`
	Service   string            `json:"service,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
}

// Logger writes audit events to an output destination.
type Logger struct {
	out io.Writer
	now func() time.Time
}

// New returns a Logger writing to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{out: w, now: time.Now}
}

// Log writes an audit event at the given level.
func (l *Logger) Log(level Level, message, service string, fields map[string]string) error {
	e := Event{
		Timestamp: l.now().UTC(),
		Level:     level,
		Message:   message,
		Service:   service,
		Fields:    fields,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintf(l.out, "%s\n", b)
	return err
}

// Info logs an informational audit event.
func (l *Logger) Info(message, service string, fields map[string]string) error {
	return l.Log(LevelInfo, message, service, fields)
}

// Warn logs a warning audit event.
func (l *Logger) Warn(message, service string, fields map[string]string) error {
	return l.Log(LevelWarn, message, service, fields)
}

// Error logs an error audit event.
func (l *Logger) Error(message, service string, fields map[string]string) error {
	return l.Log(LevelError, message, service, fields)
}
