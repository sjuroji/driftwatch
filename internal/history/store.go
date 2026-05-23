// Package history provides a simple file-backed store for persisting
// drift detection results across runs, enabling trend analysis and
// change detection over time.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Entry represents a single persisted drift detection run.
type Entry struct {
	Timestamp time.Time   `json:"timestamp"`
	Report    drift.Report `json:"report"`
}

// Store persists drift report entries to a directory on disk.
type Store struct {
	dir string
}

// NewStore creates a Store that writes entries under dir.
// The directory is created if it does not exist.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("history: create store dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes a new entry for the given report using the current UTC time.
func (s *Store) Save(report drift.Report) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Report:    report,
	}

	fileName := entry.Timestamp.Format("20060102T150405Z") + ".json"
	path := filepath.Join(s.dir, fileName)

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal entry: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("history: write entry %s: %w", path, err)
	}
	return nil
}

// List returns all stored entries sorted oldest-first.
func (s *Store) List() ([]Entry, error) {
	glob := filepath.Join(s.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("history: glob entries: %w", err)
	}

	entries := make([]Entry, 0, len(matches))
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("history: read entry %s: %w", path, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: parse entry %s: %w", path, err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
