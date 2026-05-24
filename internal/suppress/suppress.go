// Package suppress provides a mechanism for suppressing drift alerts
// for known, accepted deviations from declared manifests.
package suppress

import (
	"fmt"
	"strings"
	"time"
)

// Rule defines a suppression rule for a specific service and optional field.
type Rule struct {
	Service   string    `json:"service"`
	Field     string    `json:"field,omitempty"` // empty means suppress all drift for service
	Reason    string    `json:"reason"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired reports whether the rule has passed its expiry time.
func (r Rule) IsExpired(now time.Time) bool {
	return !r.ExpiresAt.IsZero() && now.After(r.ExpiresAt)
}

// Set holds a collection of suppression rules.
type Set struct {
	rules []Rule
	now   func() time.Time
}

// New returns a new Set. If now is nil, time.Now is used.
func New(rules []Rule, now func() time.Time) *Set {
	if now == nil {
		now = time.Now
	}
	return &Set{rules: rules, now: now}
}

// IsSuppressed reports whether drift for the given service and field is
// currently suppressed. An empty field matches any field-level suppression
// as well as service-wide suppressions.
func (s *Set) IsSuppressed(service, field string) bool {
	now := s.now()
	for _, r := range s.rules {
		if r.IsExpired(now) {
			continue
		}
		if !strings.EqualFold(r.Service, service) {
			continue
		}
		// Service-wide suppression (no specific field).
		if r.Field == "" {
			return true
		}
		if strings.EqualFold(r.Field, field) {
			return true
		}
	}
	return false
}

// Active returns all non-expired rules.
func (s *Set) Active() []Rule {
	now := s.now()
	out := make([]Rule, 0, len(s.rules))
	for _, r := range s.rules {
		if !r.IsExpired(now) {
			out = append(out, r)
		}
	}
	return out
}

// Validate checks that a Rule is well-formed.
func Validate(r Rule) error {
	if strings.TrimSpace(r.Service) == "" {
		return fmt.Errorf("suppress: rule must specify a service")
	}
	if strings.TrimSpace(r.Reason) == "" {
		return fmt.Errorf("suppress: rule for %q must include a reason", r.Service)
	}
	return nil
}
