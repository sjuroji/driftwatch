// Package tag provides utilities for matching and filtering drift reports
// based on user-defined string tags attached to manifests and live states.
package tag

import "strings"

// Set represents a collection of tags as a map for O(1) lookup.
type Set map[string]struct{}

// NewSet creates a Set from a slice of tag strings. Tags are normalised to
// lowercase so matching is case-insensitive.
func NewSet(tags []string) Set {
	s := make(Set, len(tags))
	for _, t := range tags {
		s[strings.ToLower(strings.TrimSpace(t))] = struct{}{}
	}
	return s
}

// Contains reports whether the set contains the given tag (case-insensitive).
func (s Set) Contains(tag string) bool {
	_, ok := s[strings.ToLower(strings.TrimSpace(tag))]
	return ok
}

// Intersects reports whether s and other share at least one common tag.
func (s Set) Intersects(other Set) bool {
	for t := range other {
		if _, ok := s[t]; ok {
			return true
		}
	}
	return false
}

// Slice returns the tags in the set as a sorted slice.
func (s Set) Slice() []string {
	out := make([]string, 0, len(s))
	for t := range s {
		out = append(out, t)
	}
	sortStrings(out)
	return out
}

// Match returns the subset of candidates whose tag sets intersect with
// required. If required is empty, all candidates are returned unchanged.
func Match(required Set, candidates []Tagged) []Tagged {
	if len(required) == 0 {
		return candidates
	}
	var out []Tagged
	for _, c := range candidates {
		if required.Intersects(NewSet(c.Tags())) {
			out = append(out, c)
		}
	}
	return out
}

// Tagged is implemented by any type that exposes a list of string tags.
type Tagged interface {
	Tags() []string
}

// sortStrings is a simple insertion sort to avoid importing "sort" for a
// small slice.
func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
