// Package tag provides case-insensitive tag sets and matching utilities for
// driftwatch manifests and live-state objects.
//
// Tags are arbitrary string labels (e.g. "prod", "eu-west", "critical") that
// can be attached to services. The Match function lets callers filter any
// slice of Tagged values down to only those whose tags overlap with a required
// Set, enabling tag-scoped drift detection runs.
//
// # Basic usage
//
//	required := tag.NewSet([]string{"prod"})
//	filtered := tag.Match(required, allManifests)
//
// Tags are normalised to lowercase on construction so comparisons are always
// case-insensitive without additional effort from callers.
package tag
