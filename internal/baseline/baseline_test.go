package baseline_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/driftwatch/internal/baseline"
	"github.com/driftwatch/internal/drift"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func sampleState() drift.LiveState {
	return drift.LiveState{
		Service:  "auth-service",
		Replicas: 3,
		Image:    "auth:v1.2.0",
		Env:      map[string]string{"LOG_LEVEL": "info"},
	}
}

func TestNewStore_CreatesDirectory(t *testing.T) {
	dir := t.TempDir() + "/baselines"
	_, err := baseline.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestSave_And_Load_RoundTrip(t *testing.T) {
	store, err := baseline.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	state := sampleState()
	if err := store.Save("auth-service", state, fixedTime); err != nil {
		t.Fatalf("Save: %v", err)
	}
	entry, err := store.Load("auth-service")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entry.Service != "auth-service" {
		t.Errorf("Service: got %q, want %q", entry.Service, "auth-service")
	}
	if !entry.SavedAt.Equal(fixedTime) {
		t.Errorf("SavedAt: got %v, want %v", entry.SavedAt, fixedTime)
	}
	if entry.LiveState.Image != state.Image {
		t.Errorf("Image: got %q, want %q", entry.LiveState.Image, state.Image)
	}
	if entry.LiveState.Replicas != state.Replicas {
		t.Errorf("Replicas: got %d, want %d", entry.LiveState.Replicas, state.Replicas)
	}
}

func TestLoad_NotFound(t *testing.T) {
	store, err := baseline.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	_, err = store.Load("nonexistent-service")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got: %v", err)
	}
}

func TestDelete_RemovesBaseline(t *testing.T) {
	store, err := baseline.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if err := store.Save("auth-service", sampleState(), fixedTime); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := store.Delete("auth-service"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err = store.Load("auth-service")
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected ErrNotExist after delete, got: %v", err)
	}
}

func TestDelete_Idempotent(t *testing.T) {
	store, err := baseline.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	// Deleting a non-existent baseline should not error.
	if err := store.Delete("ghost-service"); err != nil {
		t.Errorf("Delete non-existent: unexpected error: %v", err)
	}
}
