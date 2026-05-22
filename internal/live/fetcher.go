package live

import (
	"fmt"

	"github.com/driftwatch/internal/manifest"
)

// LiveState represents the actual running state of a deployed service.
type LiveState struct {
	Name     string
	Image    string
	Replicas int
	Env      map[string]string
}

// Fetcher retrieves live service state from a backend (e.g. Kubernetes, Docker).
type Fetcher interface {
	Fetch(name string) (*LiveState, error)
}

// StaticFetcher is a test/stub fetcher backed by a static map of states.
type StaticFetcher struct {
	States map[string]*LiveState
}

// Fetch returns the live state for the named service, or an error if not found.
func (f *StaticFetcher) Fetch(name string) (*LiveState, error) {
	state, ok := f.States[name]
	if !ok {
		return nil, fmt.Errorf("live: service %q not found", name)
	}
	return state, nil
}

// FetchAll retrieves live states for all manifests using the given fetcher.
// It returns a map keyed by service name and a slice of errors for any failures.
func FetchAll(fetcher Fetcher, manifests []*manifest.Manifest) (map[string]*LiveState, []error) {
	states := make(map[string]*LiveState, len(manifests))
	var errs []error

	for _, m := range manifests {
		state, err := fetcher.Fetch(m.Name)
		if err != nil {
			errs = append(errs, fmt.Errorf("fetch %q: %w", m.Name, err))
			continue
		}
		states[m.Name] = state
	}

	return states, errs
}
