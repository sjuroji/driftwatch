// Package retry provides a small, context-aware retry helper with
// exponential back-off for use across driftwatch subsystems.
//
// Basic usage:
//
//	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
//		return callExternalService()
//	})
//
// Permanent errors (errors that should not be retried) can be signalled
// by wrapping them with retry.PermFail:
//
//	return retry.PermFail(fmt.Errorf("invalid credentials: %w", err))
//
// The retry loop will stop immediately upon receiving a PermFail-wrapped
// error, propagating the underlying cause through the error chain so
// callers can still use errors.Is / errors.As.
package retry
