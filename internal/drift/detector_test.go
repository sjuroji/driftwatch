package drift_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/manifest"
)

func baseManifest() manifest.Manifest {
	return manifest.Manifest{
		Name:     "auth-service",
		Replicas: 3,
		Image:    "auth:v1.2.0",
		Port:     8080,
	}
}

func baseLiveState() drift.ServiceState {
	return drift.ServiceState{
		Name:     "auth-service",
		Replicas: 3,
		Image:    "auth:v1.2.0",
		Port:     8080,
	}
}

func TestDetect_InSync(t *testing.T) {
	result := drift.Detect(baseManifest(), baseLiveState())

	if result.Status != drift.StatusInSync {
		t.Errorf("expected status %q, got %q", drift.StatusInSync, result.Status)
	}
	if len(result.Diffs) != 0 {
		t.Errorf("expected no diffs, got %v", result.Diffs)
	}
}

func TestDetect_ReplicaDrift(t *testing.T) {
	live := baseLiveState()
	live.Replicas = 1

	result := drift.Detect(baseManifest(), live)

	if result.Status != drift.StatusDrifted {
		t.Errorf("expected status %q, got %q", drift.StatusDrifted, result.Status)
	}
	if len(result.Diffs) != 1 {
		t.Errorf("expected 1 diff, got %d: %v", len(result.Diffs), result.Diffs)
	}
}

func TestDetect_ImageDrift(t *testing.T) {
	live := baseLiveState()
	live.Image = "auth:v1.1.0"

	result := drift.Detect(baseManifest(), live)

	if result.Status != drift.StatusDrifted {
		t.Errorf("expected status %q, got %q", drift.StatusDrifted, result.Status)
	}
}

func TestDetect_MultipleDrifts(t *testing.T) {
	live := baseLiveState()
	live.Replicas = 5
	live.Port = 9090

	result := drift.Detect(baseManifest(), live)

	if result.Status != drift.StatusDrifted {
		t.Errorf("expected status %q, got %q", drift.StatusDrifted, result.Status)
	}
	if len(result.Diffs) != 2 {
		t.Errorf("expected 2 diffs, got %d: %v", len(result.Diffs), result.Diffs)
	}
}

func TestDetect_ServiceName(t *testing.T) {
	result := drift.Detect(baseManifest(), baseLiveState())

	if result.ServiceName != "auth-service" {
		t.Errorf("expected service name %q, got %q", "auth-service", result.ServiceName)
	}
}
