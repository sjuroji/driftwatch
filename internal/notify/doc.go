// Package notify provides notification mechanisms for driftwatch.
//
// When a drift detection run completes, a Notifier can be used to alert
// operators via various channels. Two implementations are provided:
//
//   - stdNotifier: writes a single summary line to an io.Writer (e.g. stderr).
//   - webhookNotifier: POSTs a JSON payload to an HTTP endpoint.
//
// Both implementations respect a Level setting that controls when
// notifications fire:
//
//   - LevelNone  — never notify
//   - LevelDrift — notify only when drift is detected (default)
//   - LevelAll   — notify after every run regardless of drift
//
// Example usage:
//
//	n := notify.New(notify.DefaultConfig())
//	if err := n.Notify(report); err != nil {
//		log.Printf("notification failed: %v", err)
//	}
package notify
