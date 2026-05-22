package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempManifest(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp manifest: %v", err)
	}
	return p
}

func TestRun_MissingManifestFlag(t *testing.T) {
	// run() returns an error when --manifest is not provided.
	// We simulate this by directly checking the validation logic.
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"driftwatch"}

	// Re-initialise flags to avoid "flag redefined" panics in test runs.
	// We call run indirectly via a helper that resets the flag set.
	err := runWithArgs([]string{})
	if err == nil {
		t.Fatal("expected error when --manifest is missing, got nil")
	}
}

func TestRun_InvalidManifestPath(t *testing.T) {
	err := runWithArgs([]string{"--manifest", "/nonexistent/path/manifest.yaml"})
	if err == nil {
		t.Fatal("expected error for missing manifest file, got nil")
	}
}

func TestRun_ValidManifest_TextOutput(t *testing.T) {
	path := writeTempManifest(t, `
name: test-service
image: nginx:1.25
replicas: 2
`)
	err := runWithArgs([]string{"--manifest", path, "--format", "text"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ValidManifest_JSONOutput(t *testing.T) {
	path := writeTempManifest(t, `
name: test-service
image: nginx:1.25
replicas: 1
`)
	err := runWithArgs([]string{"--manifest", path, "--format", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_UnknownFormat(t *testing.T) {
	path := writeTempManifest(t, `
name: test-service
image: nginx:1.25
replicas: 1
`)
	err := runWithArgs([]string{"--manifest", path, "--format", "xml"})
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}
