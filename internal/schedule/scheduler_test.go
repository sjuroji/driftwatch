package schedule_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/schedule"
)

func TestDefaultConfig(t *testing.T) {
	cfg := schedule.DefaultConfig()
	if cfg.Interval != 5*time.Minute {
		t.Errorf("expected 5m interval, got %v", cfg.Interval)
	}
	if cfg.MaxErrors != 5 {
		t.Errorf("expected MaxErrors=5, got %d", cfg.MaxErrors)
	}
}

func TestScheduler_RunsJobImmediately(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	cfg := schedule.Config{Interval: 10 * time.Second, MaxErrors: 0}
	s := schedule.New(cfg, job)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	s.Run(ctx)

	if atomic.LoadInt32(&count) < 1 {
		t.Error("expected job to run at least once immediately")
	}
}

func TestScheduler_RunsOnTick(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	cfg := schedule.Config{Interval: 20 * time.Millisecond, MaxErrors: 0}
	s := schedule.New(cfg, job)

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	s.Run(ctx)

	if got := atomic.LoadInt32(&count); got < 2 {
		t.Errorf("expected at least 2 runs, got %d", got)
	}
}

func TestScheduler_StopsOnMaxErrors(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return errors.New("simulated failure")
	}

	cfg := schedule.Config{Interval: 10 * time.Millisecond, MaxErrors: 3}
	s := schedule.New(cfg, job)

	// Run with a generous timeout; scheduler should self-stop before it.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	s.Run(ctx)

	// initial run + 3 ticked errors = 4 total calls
	if got := atomic.LoadInt32(&count); got > 5 {
		t.Errorf("scheduler did not stop near MaxErrors, ran %d times", got)
	}
}

func TestScheduler_ResetsErrorCountOnSuccess(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		n := atomic.AddInt32(&count, 1)
		if n%2 == 0 {
			return errors.New("even run fails")
		}
		return nil
	}

	cfg := schedule.Config{Interval: 15 * time.Millisecond, MaxErrors: 2}
	s := schedule.New(cfg, job)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	s.Run(ctx) // should not stop early because successes reset counter

	if got := atomic.LoadInt32(&count); got < 4 {
		t.Errorf("expected at least 4 runs with interleaved errors, got %d", got)
	}
}
