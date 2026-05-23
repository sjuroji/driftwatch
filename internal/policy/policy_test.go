package policy_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/policy"
)

func makeReport(statuses ...drift.Status) drift.Report {
	entries := make([]drift.Entry, len(statuses))
	for i, s := range statuses {
		entries[i] = drift.Entry{
			Service: fmt.Sprintf("svc-%d", i),
			Status:  s,
		}
	}
	return drift.Report{Entries: entries, GeneratedAt: time.Now()}
}

func TestEvaluate_EmptyReport_Passes(t *testing.T) {
	eval := policy.New([]policy.Rule{
		{Name: "zero-drift", MaxDriftPercent: 0, Severity: policy.SeverityCrit},
	})
	result := eval.Evaluate(drift.Report{})
	if !result.Passed {
		t.Fatal("expected pass on empty report")
	}
	if len(result.Violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(result.Violations))
	}
}

func TestEvaluate_AllInSync_Passes(t *testing.T) {
	eval := policy.New([]policy.Rule{
		{Name: "zero-drift", MaxDriftPercent: 0, Severity: policy.SeverityCrit},
	})
	report := makeReport(drift.StatusSync, drift.StatusSync, drift.StatusSync)
	result := eval.Evaluate(report)
	if !result.Passed {
		t.Fatalf("expected pass, got violations: %v", result.Violations)
	}
}

func TestEvaluate_DriftBelowThreshold_Passes(t *testing.T) {
	eval := policy.New([]policy.Rule{
		{Name: "allow-25pct", MaxDriftPercent: 25.0, Severity: policy.SeverityWarn},
	})
	// 1 of 4 = 25% — not strictly greater than threshold, should pass
	report := makeReport(drift.StatusDrift, drift.StatusSync, drift.StatusSync, drift.StatusSync)
	result := eval.Evaluate(report)
	if !result.Passed {
		t.Fatalf("expected pass at boundary, got violations: %v", result.Violations)
	}
}

func TestEvaluate_DriftExceedsThreshold_Fails(t *testing.T) {
	eval := policy.New([]policy.Rule{
		{Name: "allow-25pct", MaxDriftPercent: 25.0, Severity: policy.SeverityWarn},
	})
	// 2 of 4 = 50% — exceeds 25%
	report := makeReport(drift.StatusDrift, drift.StatusDrift, drift.StatusSync, drift.StatusSync)
	result := eval.Evaluate(report)
	if result.Passed {
		t.Fatal("expected failure when drift exceeds threshold")
	}
	if len(result.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(result.Violations))
	}
	v := result.Violations[0]
	if v.Rule.Name != "allow-25pct" {
		t.Errorf("unexpected rule name %q", v.Rule.Name)
	}
	if v.Actual != 50.0 {
		t.Errorf("expected actual 50.0, got %.1f", v.Actual)
	}
}

func TestEvaluate_MultipleRules_MultipleViolations(t *testing.T) {
	eval := policy.New([]policy.Rule{
		{Name: "strict", MaxDriftPercent: 0, Severity: policy.SeverityCrit},
		{Name: "loose", MaxDriftPercent: 40, Severity: policy.SeverityWarn},
	})
	// 3 of 4 = 75%: violates both rules
	report := makeReport(drift.StatusDrift, drift.StatusDrift, drift.StatusDrift, drift.StatusSync)
	result := eval.Evaluate(report)
	if result.Passed {
		t.Fatal("expected failure")
	}
	if len(result.Violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(result.Violations))
	}
}
