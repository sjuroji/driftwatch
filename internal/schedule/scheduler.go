// Package schedule provides periodic drift-check scheduling.
package schedule

import (
	"context"
	"log"
	"time"
)

// Job is a function that performs a single drift-check run.
type Job func(ctx context.Context) error

// Config holds scheduler configuration.
type Config struct {
	// Interval between successive drift checks.
	Interval time.Duration
	// MaxErrors is the number of consecutive errors before the scheduler stops.
	// Zero means never stop on errors.
	MaxErrors int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:  5 * time.Minute,
		MaxErrors: 5,
	}
}

// Scheduler runs a Job on a fixed interval.
type Scheduler struct {
	cfg Config
	job Job
}

// New creates a new Scheduler.
func New(cfg Config, job Job) *Scheduler {
	return &Scheduler{cfg: cfg, job: job}
}

// Run starts the scheduler and blocks until ctx is cancelled.
// It executes the job immediately, then on every tick.
func (s *Scheduler) Run(ctx context.Context) {
	if err := s.runJob(ctx); err != nil {
		log.Printf("schedule: initial run error: %v", err)
	}

	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()

	consecutiveErrors := 0
	for {
		select {
		case <-ctx.Done():
			log.Println("schedule: context cancelled, stopping scheduler")
			return
		case <-ticker.C:
			if err := s.runJob(ctx); err != nil {
				consecutiveErrors++
				log.Printf("schedule: job error (%d consecutive): %v", consecutiveErrors, err)
				if s.cfg.MaxErrors > 0 && consecutiveErrors >= s.cfg.MaxErrors {
					log.Printf("schedule: reached max consecutive errors (%d), stopping", s.cfg.MaxErrors)
					return
				}
			} else {
				consecutiveErrors = 0
			}
		}
	}
}

func (s *Scheduler) runJob(ctx context.Context) error {
	start := time.Now()
	err := s.job(ctx)
	elapsed := time.Since(start)
	if err != nil {
		log.Printf("schedule: job finished in %s with error: %v", elapsed.Round(time.Millisecond), err)
	} else {
		log.Printf("schedule: job finished in %s", elapsed.Round(time.Millisecond))
	}
	return err
}
