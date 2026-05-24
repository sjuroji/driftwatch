package suppress_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/suppress"
)

var (
	fixedNow  = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	futureExp = fixedNow.Add(24 * time.Hour)
	pastExp   = fixedNow.Add(-1 * time.Hour)
)

func fixedClock() func() time.Time {
	return func() time.Time { return fixedNow }
}

func TestIsSuppressed_ServiceWide(t *testing.T) {
	rules := []suppress.Rule{
		{Service: "auth-service", Reason: "planned maintenance", ExpiresAt: futureExp},
	}
	s := suppress.New(rules, fixedClock())

	if !s.IsSuppressed("auth-service", "replicas") {
		t.Error("expected auth-service to be suppressed for any field")
	}
	if !s.IsSuppressed("auth-service", "image") {
		t.Error("expected auth-service to be suppressed for image field")
	}
}

func TestIsSuppressed_FieldSpecific(t *testing.T) {
	rules := []suppress.Rule{
		{Service: "api-gateway", Field: "replicas", Reason: "autoscaler", ExpiresAt: futureExp},
	}
	s := suppress.New(rules, fixedClock())

	if !s.IsSuppressed("api-gateway", "replicas") {
		t.Error("expected replicas field to be suppressed")
	}
	if s.IsSuppressed("api-gateway", "image") {
		t.Error("image field should not be suppressed")
	}
}

func TestIsSuppressed_ExpiredRule(t *testing.T) {
	rules := []suppress.Rule{
		{Service: "auth-service", Reason: "old window", ExpiresAt: pastExp},
	}
	s := suppress.New(rules, fixedClock())

	if s.IsSuppressed("auth-service", "replicas") {
		t.Error("expired rule should not suppress drift")
	}
}

func TestIsSuppressed_CaseInsensitive(t *testing.T) {
	rules := []suppress.Rule{
		{Service: "Auth-Service", Field: "Replicas", Reason: "test", ExpiresAt: futureExp},
	}
	s := suppress.New(rules, fixedClock())

	if !s.IsSuppressed("auth-service", "replicas") {
		t.Error("matching should be case-insensitive")
	}
}

func TestActive_FiltersExpired(t *testing.T) {
	rules := []suppress.Rule{
		{Service: "svc-a", Reason: "active", ExpiresAt: futureExp},
		{Service: "svc-b", Reason: "expired", ExpiresAt: pastExp},
	}
	s := suppress.New(rules, fixedClock())
	active := s.Active()

	if len(active) != 1 {
		t.Fatalf("expected 1 active rule, got %d", len(active))
	}
	if active[0].Service != "svc-a" {
		t.Errorf("unexpected active rule service: %s", active[0].Service)
	}
}

func TestValidate_MissingService(t *testing.T) {
	err := suppress.Validate(suppress.Rule{Reason: "reason"})
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestValidate_MissingReason(t *testing.T) {
	err := suppress.Validate(suppress.Rule{Service: "svc"})
	if err == nil {
		t.Error("expected error for missing reason")
	}
}

func TestValidate_Valid(t *testing.T) {
	err := suppress.Validate(suppress.Rule{Service: "svc", Reason: "ok"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
