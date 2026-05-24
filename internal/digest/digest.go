// Package digest computes and compares deterministic hashes of service
// manifests and live states, enabling fast change detection without a
// full field-by-field diff.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Result holds the computed digest for a single service.
type Result struct {
	Service string `json:"service"`
	Hash    string `json:"hash"`
}

// Changed returns true when two Results for the same service have different
// hashes, indicating that the underlying state has changed.
func Changed(a, b Result) bool {
	return a.Hash != b.Hash
}

// Compute serialises v to canonical JSON and returns the SHA-256 hex digest.
// Any value that can be marshalled to JSON is accepted.
func Compute(service string, v any) (Result, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return Result{}, fmt.Errorf("digest: marshal %q: %w", service, err)
	}

	sum := sha256.Sum256(data)
	return Result{
		Service: service,
		Hash:    hex.EncodeToString(sum[:]),
	}, nil
}

// ComputeAll calls Compute for each entry in the map and returns a slice of
// Results in an unspecified order.  The first marshalling error is returned
// immediately.
func ComputeAll(states map[string]any) ([]Result, error) {
	results := make([]Result, 0, len(states))
	for svc, state := range states {
		r, err := Compute(svc, state)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

// Index converts a slice of Results into a map keyed by service name for
// O(1) look-ups.
func Index(results []Result) map[string]Result {
	m := make(map[string]Result, len(results))
	for _, r := range results {
		m[r.Service] = r
	}
	return m
}
