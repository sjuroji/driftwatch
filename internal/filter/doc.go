// Package filter provides manifest filtering utilities for driftwatch.
//
// Filters can be composed from two independent axes:
//
//	1. Name list  — include only the named services.
//	2. Label map  — include only manifests whose labels contain every
//	               specified key/value pair.
//
// When both axes are specified a manifest must satisfy both to be included.
// An empty Options struct passes all manifests through unchanged, making
// filter.Apply safe to call unconditionally in the main pipeline.
//
// Example:
//
//	filtered := filter.Apply(manifests, filter.Options{
//	    Names:  []string{"auth-service", "payment-service"},
//	    Labels: map[string]string{"env": "prod"},
//	})
package filter
