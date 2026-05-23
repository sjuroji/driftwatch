// Package schedule provides a simple interval-based scheduler for running
// periodic drift-check jobs in driftwatch.
//
// Basic usage:
//
//	cfg := schedule.DefaultConfig()
//	cfg.Interval = 10 * time.Minute
//
//	s := schedule.New(cfg, func(ctx context.Context) error {
//		return run(ctx, manifestPath, outputFormat)
//	})
//
//	s.Run(ctx) // blocks until ctx is cancelled or MaxErrors reached
//
// The job is executed once immediately on startup and then on every Interval
// tick. If MaxErrors consecutive failures occur the scheduler stops itself
// to avoid flooding downstream systems with broken requests.
package schedule
