// Package policy evaluates drift reports against configurable rules
// and produces a pass/fail verdict with per-rule details.
package policy

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// Severity represents how seriously a policy violation should be treated.
type Severity string

const (
	SeverityWarn Severity = "warn"
	SeverityCrit Severity = "crit"
)

// Rule defines a single policy constraint.
type Rule struct {
	// Name is a human-readable identifier for the rule.
	Name string
	// MaxDriftPercent is the maximum allowed percentage of drifted services.
	// A value of 0 means zero drift is tolerated.
	MaxDriftPercent float64
	// Severity controls whether a violation is a warning or critical failure.
	Severity Severity
}

// Violation records a rule that was not satisfied.
type Violation struct {
	Rule    Rule
	Actual  float64
	Message string
}

// Result is the outcome of evaluating a report against a policy set.
type Result struct {
	Passed     bool
	Violations []Violation
}

// Evaluator holds a set of rules and evaluates drift reports against them.
type Evaluator struct {
	rules []Rule
}

// New creates an Evaluator with the provided rules.
func New(rules []Rule) *Evaluator {
	return &Evaluator{rules: rules}
}

// Evaluate checks the report against every rule and returns a Result.
func (e *Evaluator) Evaluate(report drift.Report) Result {
	total := len(report.Entries)
	if total == 0 {
		return Result{Passed: true}
	}

	drifted := 0
	for _, entry := range report.Entries {
		if entry.Status == drift.StatusDrift {
			drifted++
		}
	}
	actualPct := float64(drifted) / float64(total) * 100.0

	var violations []Violation
	for _, rule := range e.rules {
		if actualPct > rule.MaxDriftPercent {
			violations = append(violations, Violation{
				Rule:   rule,
				Actual: actualPct,
				Message: fmt.Sprintf(
					"rule %q: drift %.1f%% exceeds allowed %.1f%%",
					rule.Name, actualPct, rule.MaxDriftPercent,
				),
			})
		}
	}
	return Result{
		Passed:     len(violations) == 0,
		Violations: violations,
	}
}
