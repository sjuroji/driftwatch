package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/manifest"
)

func writeTempManifest(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp manifest: %v", err)
	}
	return path
}

func TestLoad_ValidManifest(t *testing.T) {
	content := `
name: auth-service
version: "1.2.3"
image: ghcr.io/yourorg/auth-service:1.2.3
replicas: 2
environment:
  LOG_LEVEL: info
  PORT: "8080"
ports:
  - 8080
`
	path := writeTempManifest(t, content)
	m, err := manifest.Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if m.Name != "auth-service" {
		t.Errorf("expected name %q, got %q", "auth-service", m.Name)
	}
	if m.Replicas != 2 {
		t.Errorf("expected replicas 2, got %d", m.Replicas)
	}
	if m.Environment["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", m.Environment["LOG_LEVEL"])
	}
}

func TestLoad_MissingName(t *testing.T) {
	content := `image: ghcr.io/yourorg/auth-service:1.2.3
replicas: 1
`
	path := writeTempManifest(t, content)
	_, err := manifest.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing name")
	}
}

func TestLoad_NegativeReplicas(t *testing.T) {
	content := `name: bad-service
image: ghcr.io/yourorg/bad-service:latest
replicas: -1
`
	path := writeTempManifest(t, content)
	_, err := manifest.Load(path)
	if err == nil {
		t.Fatal("expected validation error for negative replicas")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := manifest.Load("/nonexistent/path/manifest.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
