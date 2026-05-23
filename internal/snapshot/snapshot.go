// Package snapshot captures and compares point-in-time live state
// for a set of services, enabling drift trend analysis over time.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry holds the live state of a single service at a moment in time.
type Entry struct {
	Service  string            `json:"service"`
	Image    string            `json:"image"`
	Replicas int               `json:"replicas"`
	Labels   map[string]string `json:"labels,omitempty"`
}

// Snapshot represents the full captured state across all services.
type Snapshot struct {
	CapturedAt time.Time `json:"captured_at"`
	Entries    []Entry   `json:"entries"`
}

// Store persists and retrieves snapshots from a directory on disk.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create store dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes a snapshot to disk using the captured timestamp as the filename.
func (s *Store) Save(snap Snapshot) error {
	filename := snap.CapturedAt.UTC().Format("20060102T150405Z") + ".json"
	path := filepath.Join(s.dir, filename)

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write file: %w", err)
	}
	return nil
}

// Latest returns the most recently saved snapshot, or an error if none exist.
func (s *Store) Latest() (Snapshot, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read dir: %w", err)
	}

	// ReadDir returns entries sorted by name; last entry is newest timestamp.
	var last os.DirEntry
	for i := len(entries) - 1; i >= 0; i-- {
		if !entries[i].IsDir() {
			last = entries[i]
			break
		}
	}
	if last == nil {
		return Snapshot{}, fmt.Errorf("snapshot: no snapshots found in %s", s.dir)
	}

	data, err := os.ReadFile(filepath.Join(s.dir, last.Name()))
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read file: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return snap, nil
}
