package filter_test

import (
	"testing"

	"github.com/driftwatch/internal/filter"
	"github.com/driftwatch/internal/manifest"
)

func makeManifests() []manifest.Manifest {
	return []manifest.Manifest{
		{Name: "auth-service", Image: "auth:1.0", Replicas: 2, Labels: map[string]string{"env": "prod", "team": "platform"}},
		{Name: "payment-service", Image: "payment:2.1", Replicas: 3, Labels: map[string]string{"env": "prod", "team": "payments"}},
		{Name: "notification-service", Image: "notify:0.9", Replicas: 1, Labels: map[string]string{"env": "staging", "team": "platform"}},
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	result := filter.Apply(makeManifests(), filter.Options{})
	if len(result) != 3 {
		t.Fatalf("expected 3 manifests, got %d", len(result))
	}
}

func TestApply_FilterByName(t *testing.T) {
	opts := filter.Options{Names: []string{"auth-service"}}
	result := filter.Apply(makeManifests(), opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(result))
	}
	if result[0].Name != "auth-service" {
		t.Errorf("unexpected name %q", result[0].Name)
	}
}

func TestApply_FilterByName_CaseInsensitive(t *testing.T) {
	opts := filter.Options{Names: []string{"AUTH-SERVICE"}}
	result := filter.Apply(makeManifests(), opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(result))
	}
}

func TestApply_FilterByLabel(t *testing.T) {
	opts := filter.Options{Labels: map[string]string{"env": "prod"}}
	result := filter.Apply(makeManifests(), opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 manifests, got %d", len(result))
	}
}

func TestApply_FilterByMultipleLabels(t *testing.T) {
	opts := filter.Options{Labels: map[string]string{"env": "prod", "team": "platform"}}
	result := filter.Apply(makeManifests(), opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(result))
	}
	if result[0].Name != "auth-service" {
		t.Errorf("unexpected name %q", result[0].Name)
	}
}

func TestApply_FilterByNameAndLabel(t *testing.T) {
	opts := filter.Options{
		Names:  []string{"payment-service"},
		Labels: map[string]string{"env": "prod"},
	}
	result := filter.Apply(makeManifests(), opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(result))
	}
}

func TestApply_NoMatch_ReturnsEmpty(t *testing.T) {
	opts := filter.Options{Names: []string{"unknown-service"}}
	result := filter.Apply(makeManifests(), opts)
	if len(result) != 0 {
		t.Fatalf("expected 0 manifests, got %d", len(result))
	}
}

func TestApply_EmptyInput(t *testing.T) {
	result := filter.Apply(nil, filter.Options{Names: []string{"auth-service"}})
	if result != nil && len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}
