package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "driftwatch.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return path
}

func TestLoad_EmptyPath_ReturnsDefaults(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("expected default output_format=text, got %q", cfg.OutputFormat)
	}
	if cfg.ManifestDir != "./manifests" {
		t.Errorf("expected default manifest_dir=./manifests, got %q", cfg.ManifestDir)
	}
	if cfg.FailOnDrift {
		t.Error("expected fail_on_drift=false by default")
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
manifest_dir: /etc/manifests
output_format: json
fail_on_drift: true
services:
  - auth-service
  - billing-service
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ManifestDir != "/etc/manifests" {
		t.Errorf("manifest_dir: got %q", cfg.ManifestDir)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("output_format: got %q", cfg.OutputFormat)
	}
	if !cfg.FailOnDrift {
		t.Error("expected fail_on_drift=true")
	}
	if len(cfg.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(cfg.Services))
	}
}

func TestLoad_InvalidOutputFormat(t *testing.T) {
	path := writeTempConfig(t, `output_format: xml`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid output_format, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_MalformedYAML(t *testing.T) {
	path := writeTempConfig(t, `output_format: [unclosed`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for malformed YAML, got nil")
	}
}
