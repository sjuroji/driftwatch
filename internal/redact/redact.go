// Package redact provides utilities for scrubbing sensitive values
// from manifest fields and live state before logging or reporting.
package redact

import "strings"

// DefaultSensitiveKeys is the list of field keys considered sensitive
// by default. Matching is case-insensitive.
var DefaultSensitiveKeys = []string{
	"password",
	"secret",
	"token",
	"api_key",
	"apikey",
	"private_key",
	"privatekey",
	"auth",
	"credential",
}

const redactedValue = "[REDACTED]"

// Options controls which keys are treated as sensitive.
type Options struct {
	// SensitiveKeys overrides DefaultSensitiveKeys when non-nil.
	SensitiveKeys []string
}

// Map returns a copy of m with sensitive values replaced by [REDACTED].
// Keys are matched case-insensitively against the configured sensitive list.
func Map(m map[string]string, opts Options) map[string]string {
	keys := opts.SensitiveKeys
	if keys == nil {
		keys = DefaultSensitiveKeys
	}

	out := make(map[string]string, len(m))
	for k, v := range m {
		if isSensitive(k, keys) {
			out[k] = redactedValue
		} else {
			out[k] = v
		}
	}
	return out
}

// Value returns [REDACTED] if key is sensitive, otherwise returns value unchanged.
func Value(key, value string, opts Options) string {
	keys := opts.SensitiveKeys
	if keys == nil {
		keys = DefaultSensitiveKeys
	}
	if isSensitive(key, keys) {
		return redactedValue
	}
	return value
}

func isSensitive(key string, sensitiveKeys []string) bool {
	lower := strings.ToLower(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(lower, strings.ToLower(s)) {
			return true
		}
	}
	return false
}
