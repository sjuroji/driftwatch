package drift

import (
	"fmt"

	"github.com/driftwatch/internal/manifest"
)

// Status represents the drift state of a service.
type Status string

const (
	StatusInSync  Status = "in_sync"
	StatusDrifted Status = "drifted"
	StatusUnknown Status = "unknown"
)

// ServiceState represents the live state of a deployed service.
type ServiceState struct {
	Name     string
	Replicas int
	Image    string
	Port     int
}

// Result holds the outcome of a drift check for a single service.
type Result struct {
	ServiceName string
	Status      Status
	Diffs       []string
}

// Detect compares a declared manifest against the live service state
// and returns a Result describing any detected drift.
func Detect(m manifest.Manifest, live ServiceState) Result {
	result := Result{
		ServiceName: m.Name,
		Status:      StatusInSync,
	}

	if m.Name != live.Name {
		result.Diffs = append(result.Diffs,
			fmt.Sprintf("name: declared=%q live=%q", m.Name, live.Name))
	}

	if m.Replicas != live.Replicas {
		result.Diffs = append(result.Diffs,
			fmt.Sprintf("replicas: declared=%d live=%d", m.Replicas, live.Replicas))
	}

	if m.Image != live.Image {
		result.Diffs = append(result.Diffs,
			fmt.Sprintf("image: declared=%q live=%q", m.Image, live.Image))
	}

	if m.Port != live.Port {
		result.Diffs = append(result.Diffs,
			fmt.Sprintf("port: declared=%d live=%d", m.Port, live.Port))
	}

	if len(result.Diffs) > 0 {
		result.Status = StatusDrifted
	}

	return result
}
