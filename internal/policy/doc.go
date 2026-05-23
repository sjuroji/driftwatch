// Package policy provides rule-based evaluation of drift reports.
//
// An Evaluator is configured with one or more Rules, each specifying
// a MaxDriftPercent threshold and a Severity level. When Evaluate is
// called with a drift.Report, the evaluator computes the percentage of
// drifted services and checks it against every rule, collecting
// Violations for any rule whose threshold is exceeded.
//
// Example usage:
//
//	eval := policy.New([]policy.Rule{
//		{Name: "zero-drift", MaxDriftPercent: 0, Severity: policy.SeverityCrit},
//	})
//	result := eval.Evaluate(report)
//	if !result.Passed {
//		for _, v := range result.Violations {
//			log.Println(v.Message)
//		}
//	}
package policy
