package debounce_test

import (
	"testing"
	"time"

	"driftwatch/internal/debounce"
)

var (
	epoch    = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	cooldown = 10 * time.Minute
)

func fixedClock(t time.Time) debounce.Clock {
	return func() time.Time { return t }
}

func TestAllow_FirstCallAlwaysTrue(t *testing.T) {
	d := debounce.New(cooldown, fixedClock(epoch))
	if !d.Allow("svc-a") {
		t.Fatal("expected first Allow to return true")
	}
}

func TestAllow_SecondCallWithinCooldown_False(t *testing.T) {
	current := epoch
	clock := func() time.Time { return current }
	d := debounce.New(cooldown, clock)

	d.Allow("svc-a") // first fire
	current = epoch.Add(5 * time.Minute) // still within cooldown

	if d.Allow("svc-a") {
		t.Fatal("expected Allow to return false within cooldown window")
	}
}

func TestAllow_AfterCooldownElapsed_True(t *testing.T) {
	current := epoch
	clock := func() time.Time { return current }
	d := debounce.New(cooldown, clock)

	d.Allow("svc-a")
	current = epoch.Add(11 * time.Minute) // past cooldown

	if !d.Allow("svc-a") {
		t.Fatal("expected Allow to return true after cooldown elapsed")
	}
}

func TestAllow_DifferentServices_Independent(t *testing.T) {
	current := epoch
	clock := func() time.Time { return current }
	d := debounce.New(cooldown, clock)

	d.Allow("svc-a")
	current = epoch.Add(2 * time.Minute)

	// svc-b has never fired — should be allowed
	if !d.Allow("svc-b") {
		t.Fatal("expected svc-b to be allowed independently of svc-a")
	}
	// svc-a is still in cooldown
	if d.Allow("svc-a") {
		t.Fatal("expected svc-a to still be suppressed")
	}
}

func TestReset_AllowsImmediateRefire(t *testing.T) {
	current := epoch
	clock := func() time.Time { return current }
	d := debounce.New(cooldown, clock)

	d.Allow("svc-a")
	current = epoch.Add(1 * time.Minute) // still in cooldown
	d.Reset("svc-a")

	if !d.Allow("svc-a") {
		t.Fatal("expected Allow to return true after Reset")
	}
}

func TestResetAll_ClearsAllServices(t *testing.T) {
	current := epoch
	clock := func() time.Time { return current }
	d := debounce.New(cooldown, clock)

	d.Allow("svc-a")
	d.Allow("svc-b")
	current = epoch.Add(1 * time.Minute)
	d.ResetAll()

	if !d.Allow("svc-a") || !d.Allow("svc-b") {
		t.Fatal("expected all services to be allowed after ResetAll")
	}
}

func TestNew_ZeroCooldown_UsesDefault(t *testing.T) {
	d := debounce.New(0, fixedClock(epoch))
	if d == nil {
		t.Fatal("expected non-nil Debouncer")
	}
	// Should not panic and first Allow should be true
	if !d.Allow("svc") {
		t.Fatal("expected first Allow to return true with default cooldown")
	}
}
