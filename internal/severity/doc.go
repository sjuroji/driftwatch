// Package severity classifies the overall severity of a drift report.
//
// An Evaluator counts the number of drifted services in a [drift.Report] and
// maps that count to one of five levels: none, low, medium, high, or critical.
//
// Thresholds are fully configurable via [Config]; [DefaultConfig] provides
// sensible production defaults.
//
// Typical usage:
//
//	eval := severity.New(severity.DefaultConfig())
//	level := eval.Evaluate(report)
//	fmt.Println("drift severity:", level)
package severity
