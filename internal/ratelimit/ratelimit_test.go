package ratelimit

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_InvalidRate(t *testing.T) {
	_, err := New(Config{Rate: 0, Interval: time.Minute})
	if err == nil {
		t.Fatal("expected error for Rate=0")
	}
}

func TestNew_InvalidInterval(t *testing.T) {
	_, err := New(Config{Rate: 1, Interval: 0})
	if err == nil {
		t.Fatal("expected error for Interval=0")
	}
}

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	l, _ := New(DefaultConfig())
	if !l.Allow("svc-a") {
		t.Fatal("first call should be allowed")
	}
}

func TestAllow_WithinRateLimit(t *testing.T) {
	base := time.Now()
	l, _ := New(Config{Rate: 3, Interval: time.Minute})
	l.now = fixedClock(base)

	for i := 0; i < 3; i++ {
		if !l.Allow("svc-b") {
			t.Fatalf("call %d should be allowed", i+1)
		}
	}
}

func TestAllow_ExceedsRateLimit(t *testing.T) {
	base := time.Now()
	l, _ := New(Config{Rate: 2, Interval: time.Minute})
	l.now = fixedClock(base)

	l.Allow("svc-c")
	l.Allow("svc-c")
	if l.Allow("svc-c") {
		t.Fatal("third call should be denied")
	}
}

func TestAllow_WindowResets_AfterInterval(t *testing.T) {
	base := time.Now()
	l, _ := New(Config{Rate: 1, Interval: time.Minute})
	l.now = fixedClock(base)

	l.Allow("svc-d") // consumes the single slot
	if l.Allow("svc-d") {
		t.Fatal("second call in same window should be denied")
	}

	// advance past the window
	l.now = fixedClock(base.Add(2 * time.Minute))
	if !l.Allow("svc-d") {
		t.Fatal("call after window reset should be allowed")
	}
}

func TestAllow_DifferentKeys_Independent(t *testing.T) {
	base := time.Now()
	l, _ := New(Config{Rate: 1, Interval: time.Minute})
	l.now = fixedClock(base)

	l.Allow("svc-x") // exhaust svc-x
	if !l.Allow("svc-y") {
		t.Fatal("svc-y should be independent of svc-x")
	}
}

func TestReset_ClearsState(t *testing.T) {
	base := time.Now()
	l, _ := New(Config{Rate: 1, Interval: time.Minute})
	l.now = fixedClock(base)

	l.Allow("svc-e") // exhaust
	l.Reset("svc-e")
	if !l.Allow("svc-e") {
		t.Fatal("after reset, call should be allowed again")
	}
}
