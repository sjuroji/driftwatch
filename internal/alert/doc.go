// Package alert provides threshold-based alerting for the driftwatch tool.
//
// An Alerter is configured with warn and critical thresholds representing
// the number of drifted services required to trigger each level. Call
// Evaluate with the current drifted-service count after each detection run;
// it returns a non-nil *Alert when a threshold is crossed and writes a
// human-readable line to the configured writer (defaulting to os.Stderr).
//
// Example usage:
//
//	 alerter := alert.New(alert.DefaultConfig(), os.Stderr)
//	 if a := alerter.Evaluate(driftedCount); a != nil {
//	     log.Printf("alert triggered: %s", a.Message)
//	 }
package alert
