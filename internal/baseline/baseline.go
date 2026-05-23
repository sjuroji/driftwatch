// Package baseline manages the storage and retrieval of known-good
// configuration states that drift detection compares against.
package baseline

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry represents a saved baseline for a single service.
type Entry struct {
	Service   string          `json:"service"`
	SavedAt   time.Time       `json:"saved_at"`
	LiveState drift.LiveState `json:"live_state"`
}

// Store persists baseline entries to disk as JSON files.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("baseline: create directory: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes the live state for the named service as the new baseline.
func (s *Store) Save(service string, state drift.LiveState, now time.Time) error {
	entry := Entry{
		Service:   service,
		SavedAt:   now,
		LiveState: state,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal %s: %w", service, err)
	}
	path := filepath.Join(s.dir, service+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", service, err)
	}
	return nil
}

// Load retrieves the saved baseline for the named service.
// Returns an error wrapping os.ErrNotExist if no baseline has been saved.
func (s *Store) Load(service string) (Entry, error) {
	path := filepath.Join(s.dir, service+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Entry{}, fmt.Errorf("baseline: no baseline for %q: %w", service, os.ErrNotExist)
		}
		return Entry{}, fmt.Errorf("baseline: read %s: %w", service, err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return Entry{}, fmt.Errorf("baseline: unmarshal %s: %w", service, err)
	}
	return entry, nil
}

// Delete removes the saved baseline for the named service.
func (s *Store) Delete(service string) error {
	path := filepath.Join(s.dir, service+".json")
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("baseline: delete %s: %w", service, err)
	}
	return nil
}
