// Package retry provides configurable retry logic with exponential backoff
// for transient failures encountered during drift detection runs.
package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Config holds retry behaviour settings.
type Config struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int `yaml:"max_attempts"`
	// InitialDelay is the wait time before the second attempt.
	InitialDelay time.Duration `yaml:"initial_delay"`
	// MaxDelay caps the exponential back-off ceiling.
	MaxDelay time.Duration `yaml:"max_delay"`
	// Multiplier scales the delay after each failure.
	Multiplier float64 `yaml:"multiplier"`
}

// DefaultConfig returns a Config suitable for most network calls.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:  3,
		InitialDelay: 200 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}
}

// Do calls fn up to cfg.MaxAttempts times, backing off between retries.
// It stops early if ctx is cancelled or fn returns a non-retryable error
// wrapped with PermFail.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	delay := cfg.InitialDelay
	var last error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry cancelled after %d attempt(s): %w", attempt-1, err)
		}
		last = fn()
		if last == nil {
			return nil
		}
		var pf *permError
		if errors.As(last, &pf) {
			return last
		}
		if attempt == cfg.MaxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		case <-time.After(delay):
		}
		delay = time.Duration(float64(delay) * cfg.Multiplier)
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
	}
	return fmt.Errorf("all %d attempt(s) failed: %w", cfg.MaxAttempts, last)
}

// permError marks an error as permanent (non-retryable).
type permError struct{ cause error }

func (p *permError) Error() string { return p.cause.Error() }
func (p *permError) Unwrap() error { return p.cause }

// PermFail wraps err so that Do treats it as permanent and stops retrying.
func PermFail(err error) error { return &permError{cause: err} }
