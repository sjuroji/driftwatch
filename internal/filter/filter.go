// Package filter provides utilities for selecting a subset of manifests
// based on label selectors and service name patterns.
package filter

import (
	"strings"

	"github.com/driftwatch/internal/manifest"
)

// Options holds the criteria used to filter a slice of manifests.
type Options struct {
	// Names is an optional list of service names to include.
	// An empty slice means "include all".
	Names []string

	// Labels is an optional map of key/value pairs that must all be present
	// on a manifest for it to be included.
	Labels map[string]string
}

// Apply returns the subset of manifests that match all criteria in opts.
func Apply(manifests []manifest.Manifest, opts Options) []manifest.Manifest {
	var out []manifest.Manifest
	for _, m := range manifests {
		if !matchesNames(m, opts.Names) {
			continue
		}
		if !matchesLabels(m, opts.Labels) {
			continue
		}
		out = append(out, m)
	}
	return out
}

// matchesNames returns true when the manifest name appears in the allow-list,
// or when the allow-list is empty (no restriction).
func matchesNames(m manifest.Manifest, names []string) bool {
	if len(names) == 0 {
		return true
	}
	for _, n := range names {
		if strings.EqualFold(m.Name, n) {
			return true
		}
	}
	return false
}

// matchesLabels returns true when every key/value pair in want is present in
// the manifest labels, or when want is empty.
func matchesLabels(m manifest.Manifest, want map[string]string) bool {
	if len(want) == 0 {
		return true
	}
	for k, v := range want {
		if m.Labels[k] != v {
			return false
		}
	}
	return true
}
