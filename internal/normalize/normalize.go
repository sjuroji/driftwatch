// Package normalize provides utilities for normalising service names and
// field keys before comparison, ensuring consistent drift detection across
// different naming conventions (e.g. camelCase vs kebab-case vs snake_case).
package normalize

import (
	"regexp"
	"strings"
)

var (
	// nonAlphanumeric matches any character that is not a letter or digit.
	nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

	// camelBoundary matches boundaries between a lowercase letter and an
	// uppercase letter, used to split camelCase identifiers.
	camelBoundary = regexp.MustCompile(`([a-z])([A-Z])`)
)

// ServiceName returns a canonical lower-kebab-case form of a service name.
// It handles camelCase, snake_case, and kebab-case inputs uniformly so that
// "AuthService", "auth_service", and "auth-service" all produce "auth-service".
func ServiceName(name string) string {
	if name == "" {
		return ""
	}
	// Split camelCase boundaries before lowercasing.
	name = camelBoundary.ReplaceAllString(name, "${1}-${2}")
	name = strings.ToLower(name)
	name = nonAlphanumeric.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	return name
}

// FieldKey returns a canonical lower-snake-case form of a field key.
// It normalises camelCase, kebab-case, and mixed inputs so that "ImageTag",
// "image-tag", and "image_tag" all produce "image_tag".
func FieldKey(key string) string {
	if key == "" {
		return ""
	}
	key = camelBoundary.ReplaceAllString(key, "${1}_${2}")
	key = strings.ToLower(key)
	key = nonAlphanumeric.ReplaceAllString(key, "_")
	key = strings.Trim(key, "_")
	return key
}

// Labels returns a new map with all keys normalised via FieldKey and all
// string values trimmed of leading/trailing whitespace. The original map is
// never mutated.
func Labels(labels map[string]string) map[string]string {
	out := make(map[string]string, len(labels))
	for k, v := range labels {
		out[FieldKey(k)] = strings.TrimSpace(v)
	}
	return out
}
