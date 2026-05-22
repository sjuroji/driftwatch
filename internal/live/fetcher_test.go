package live_test

import (
	"testing"

	"github.com/driftwatch/internal/live"
	"github.com/driftwatch/internal/manifest"
)

func TestStaticFetcher_Found(t *testing.T) {
	f := &live.StaticFetcher{
		States: map[string]*live.LiveState{
			"auth-service": {
				Name:     "auth-service",
				Image:    "auth:v1.2.3",
				Replicas: 3,
				Env:      map[string]string{"LOG_LEVEL": "info"},
			},
		},
	}

	state, err := f.Fetch("auth-service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Image != "auth:v1.2.3" {
		t.Errorf("expected image auth:v1.2.3, got %s", state.Image)
	}
	if state.Replicas != 3 {
		t.Errorf("expected 3 replicas, got %d", state.Replicas)
	}
}

func TestStaticFetcher_NotFound(t *testing.T) {
	f := &live.StaticFetcher{States: map[string]*live.LiveState{}}

	_, err := f.Fetch("missing-service")
	if err == nil {
		t.Fatal("expected error for missing service, got nil")
	}
}

func TestFetchAll_AllFound(t *testing.T) {
	f := &live.StaticFetcher{
		States: map[string]*live.LiveState{
			"svc-a": {Name: "svc-a", Image: "img-a:1", Replicas: 1},
			"svc-b": {Name: "svc-b", Image: "img-b:2", Replicas: 2},
		},
	}
	manifests := []*manifest.Manifest{
		{Name: "svc-a", Image: "img-a:1", Replicas: 1},
		{Name: "svc-b", Image: "img-b:2", Replicas: 2},
	}

	states, errs := live.FetchAll(f, manifests)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(states) != 2 {
		t.Errorf("expected 2 states, got %d", len(states))
	}
}

func TestFetchAll_PartialMissing(t *testing.T) {
	f := &live.StaticFetcher{
		States: map[string]*live.LiveState{
			"svc-a": {Name: "svc-a", Image: "img-a:1", Replicas: 1},
		},
	}
	manifests := []*manifest.Manifest{
		{Name: "svc-a", Image: "img-a:1", Replicas: 1},
		{Name: "svc-missing", Image: "img-x:1", Replicas: 1},
	}

	states, errs := live.FetchAll(f, manifests)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
	if len(states) != 1 {
		t.Errorf("expected 1 state, got %d", len(states))
	}
}
