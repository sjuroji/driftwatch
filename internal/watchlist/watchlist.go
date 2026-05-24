// Package watchlist manages the set of services actively monitored by driftwatch.
// It supports adding, removing, and querying services with optional tag filtering.
package watchlist

import (
	"errors"
	"strings"
	"sync"
)

// Entry represents a single watched service.
type Entry struct {
	Name   string            `json:"name"`
	Tags   []string          `json:"tags,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

// Watchlist holds the set of services under active monitoring.
type Watchlist struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an empty Watchlist.
func New() *Watchlist {
	return &Watchlist{
		entries: make(map[string]Entry),
	}
}

// Add inserts or replaces an entry in the watchlist.
// Returns an error if the entry name is empty.
func (w *Watchlist) Add(e Entry) error {
	if strings.TrimSpace(e.Name) == "" {
		return errors.New("watchlist: entry name must not be empty")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries[strings.ToLower(e.Name)] = e
	return nil
}

// Remove deletes a service from the watchlist by name.
// Returns false if the service was not present.
func (w *Watchlist) Remove(name string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	key := strings.ToLower(name)
	_, ok := w.entries[key]
	if ok {
		delete(w.entries, key)
	}
	return ok
}

// Contains reports whether a service is in the watchlist.
func (w *Watchlist) Contains(name string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	_, ok := w.entries[strings.ToLower(name)]
	return ok
}

// All returns a snapshot of all current entries.
func (w *Watchlist) All() []Entry {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]Entry, 0, len(w.entries))
	for _, e := range w.entries {
		out = append(out, e)
	}
	return out
}

// FilterByTag returns entries that carry at least one of the given tags.
// If tags is empty, all entries are returned.
func (w *Watchlist) FilterByTag(tags ...string) []Entry {
	if len(tags) == 0 {
		return w.All()
	}
	want := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		want[strings.ToLower(t)] = struct{}{}
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	var out []Entry
	for _, e := range w.entries {
		for _, t := range e.Tags {
			if _, ok := want[strings.ToLower(t)]; ok {
				out = append(out, e)
				break
			}
		}
	}
	return out
}
