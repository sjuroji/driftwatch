package groupby_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/groupby"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Status: drift.StatusInSync, Labels: map[string]string{"env": "prod"}},
		{Service: "auth", Status: drift.StatusDrifted, Labels: map[string]string{"env": "prod"}},
		{Service: "worker", Status: drift.StatusDrifted, Labels: map[string]string{"env": "staging"}},
		{Service: "cron", Status: drift.StatusInSync, Labels: map[string]string{}},
	}
}

func TestGroup_ByLabel(t *testing.T) {
	results := makeResults()
	fn := groupby.ByLabel("env", "unknown")
	groups := groupby.Group(results, fn)

	if len(groups["prod"]) != 2 {
		t.Errorf("expected 2 prod results, got %d", len(groups["prod"]))
	}
	if len(groups["staging"]) != 1 {
		t.Errorf("expected 1 staging result, got %d", len(groups["staging"]))
	}
	if len(groups["unknown"]) != 1 {
		t.Errorf("expected 1 unknown result, got %d", len(groups["unknown"]))
	}
}

func TestGroup_ByStatus(t *testing.T) {
	results := makeResults()
	groups := groupby.Group(results, groupby.ByStatus())

	if len(groups[drift.StatusInSync]) != 2 {
		t.Errorf("expected 2 in-sync, got %d", len(groups[drift.StatusInSync]))
	}
	if len(groups[drift.StatusDrifted]) != 2 {
		t.Errorf("expected 2 drifted, got %d", len(groups[drift.StatusDrifted]))
	}
}

func TestGroup_EmptyResults(t *testing.T) {
	groups := groupby.Group(nil, groupby.ByStatus())
	if len(groups) != 0 {
		t.Errorf("expected empty map, got %d keys", len(groups))
	}
}

func TestSummarise_Counts(t *testing.T) {
	results := makeResults()
	groups := groupby.Group(results, groupby.ByLabel("env", "unknown"))
	summaries := groupby.Summarise(groups)

	byKey := make(map[string]groupby.Summary)
	for _, s := range summaries {
		byKey[s.Key] = s
	}

	prod := byKey["prod"]
	if prod.Total != 2 || prod.Drifted != 1 || prod.InSync != 1 {
		t.Errorf("prod summary mismatch: %+v", prod)
	}

	staging := byKey["staging"]
	if staging.Total != 1 || staging.Drifted != 1 || staging.InSync != 0 {
		t.Errorf("staging summary mismatch: %+v", staging)
	}
}

func TestValidate_NilKeyFunc(t *testing.T) {
	if err := groupby.Validate(nil); err == nil {
		t.Error("expected error for nil KeyFunc")
	}
}

func TestValidate_ValidKeyFunc(t *testing.T) {
	if err := groupby.Validate(groupby.ByStatus()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
