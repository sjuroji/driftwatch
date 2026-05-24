package watchlist_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/watchlist"
)

func writeTempWatchlist(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "watchlist.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp watchlist: %v", err)
	}
	return p
}

const validWatchlist = `
services:
  - name: auth-service
    tags: [prod, critical]
    labels:
      env: production
  - name: billing
    tags: [prod]
`

func TestLoadFile_ValidFile(t *testing.T) {
	p := writeTempWatchlist(t, validWatchlist)
	w, err := watchlist.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !w.Contains("auth-service") {
		t.Error("expected auth-service in watchlist")
	}
	if !w.Contains("billing") {
		t.Error("expected billing in watchlist")
	}
}

func TestLoadFile_TagsPreserved(t *testing.T) {
	p := writeTempWatchlist(t, validWatchlist)
	w, err := watchlist.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	results := w.FilterByTag("critical")
	if len(results) != 1 || results[0].Name != "auth-service" {
		t.Errorf("expected auth-service for critical tag, got %v", results)
	}
}

func TestLoadFile_FileNotFound(t *testing.T) {
	_, err := watchlist.LoadFile("/nonexistent/watchlist.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFile_InvalidYAML(t *testing.T) {
	p := writeTempWatchlist(t, ": : invalid: yaml: [")
	_, err := watchlist.LoadFile(p)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadFile_EmptyNameEntry_ReturnsError(t *testing.T) {
	content := "services:\n  - name: \"\"\n"
	p := writeTempWatchlist(t, content)
	_, err := watchlist.LoadFile(p)
	if err == nil {
		t.Error("expected error for entry with empty name")
	}
}

func TestLoadFile_EmptyFile_ReturnsEmptyWatchlist(t *testing.T) {
	p := writeTempWatchlist(t, "services: []\n")
	w, err := watchlist.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(w.All()) != 0 {
		t.Errorf("expected empty watchlist, got %d entries", len(w.All()))
	}
}
