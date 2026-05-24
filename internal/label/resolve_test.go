package label_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/label"
)

type service struct {
	name   string
	labels map[string]string
}

func (s service) GetLabels() map[string]string { return s.labels }

func makeServices() []service {
	return []service{
		{name: "auth", labels: map[string]string{"env": "prod", "team": "platform"}},
		{name: "billing", labels: map[string]string{"env": "prod", "team": "finance"}},
		{name: "worker", labels: map[string]string{"env": "staging", "team": "platform"}},
	}
}

func TestFilter_EmptySelector_ReturnsAll(t *testing.T) {
	sel, _ := label.ParseSelector("")
	got := label.Filter(makeServices(), sel)
	if len(got) != 3 {
		t.Errorf("expected 3 items, got %d", len(got))
	}
}

func TestFilter_ByEnv(t *testing.T) {
	sel, _ := label.ParseSelector("env=prod")
	got := label.Filter(makeServices(), sel)
	if len(got) != 2 {
		t.Errorf("expected 2 prod services, got %d", len(got))
	}
}

func TestFilter_MultiLabel(t *testing.T) {
	sel, _ := label.ParseSelector("env=prod,team=platform")
	got := label.Filter(makeServices(), sel)
	if len(got) != 1 || got[0].name != "auth" {
		t.Errorf("expected only auth service, got %v", got)
	}
}

func TestFilter_NoMatch_ReturnsEmpty(t *testing.T) {
	sel, _ := label.ParseSelector("env=canary")
	got := label.Filter(makeServices(), sel)
	if len(got) != 0 {
		t.Errorf("expected 0 results, got %d", len(got))
	}
}

func TestGroupBy_Team(t *testing.T) {
	groups := label.GroupBy(makeServices(), "team")
	if len(groups["platform"]) != 2 {
		t.Errorf("expected 2 platform services, got %d", len(groups["platform"]))
	}
	if len(groups["finance"]) != 1 {
		t.Errorf("expected 1 finance service, got %d", len(groups["finance"]))
	}
}

func TestGroupBy_MissingLabel_EmptyKey(t *testing.T) {
	items := []service{
		{name: "orphan", labels: map[string]string{}},
	}
	groups := label.GroupBy(items, "team")
	if len(groups[""] ) != 1 {
		t.Errorf("expected orphan under empty key")
	}
}
