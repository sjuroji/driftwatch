// Package label provides utilities for matching and filtering services
// based on key/value label selectors, similar to Kubernetes label selectors.
package label

import (
	"fmt"
	"strings"
)

// Selector represents a set of required key/value label pairs.
type Selector map[string]string

// ParseSelector parses a comma-separated list of key=value pairs into a Selector.
// Example input: "env=prod,team=platform"
func ParseSelector(raw string) (Selector, error) {
	if raw == "" {
		return Selector{}, nil
	}
	sel := make(Selector)
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 || kv[0] == "" {
			return nil, fmt.Errorf("invalid selector segment %q: expected key=value", part)
		}
		sel[strings.ToLower(kv[0])] = kv[1]
	}
	return sel, nil
}

// Matches reports whether the given labels satisfy all requirements in the selector.
// Label keys are matched case-insensitively.
func (s Selector) Matches(labels map[string]string) bool {
	for k, want := range s {
		got, ok := labels[strings.ToLower(k)]
		if !ok || got != want {
			return false
		}
	}
	return true
}

// Empty reports whether the selector has no requirements.
func (s Selector) Empty() bool {
	return len(s) == 0
}

// String returns a canonical string representation of the selector.
func (s Selector) String() string {
	if s.Empty() {
		return ""
	}
	parts := make([]string, 0, len(s))
	for k, v := range s {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, ",")
}
