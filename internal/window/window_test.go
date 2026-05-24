package window

import (
	"testing"
	"time"
)

// fixedClock returns a function that advances by delta on each call.
func fixedClock(start time.Time, delta time.Duration) func() time.Time {
	t := start
	return func() time.Time {
		current := t
		t = t.Add(delta)
		return current
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Size != 5*time.Minute {
		t.Errorf("expected 5m, got %s", cfg.Size)
	}
}

func TestNew_InvalidSize(t *testing.T) {
	_, err := New(Config{Size: 0})
	if err == nil {
		t.Fatal("expected error for zero size")
	}
}

func TestCount_EmptyCounter(t *testing.T) {
	c, _ := New(DefaultConfig())
	if got := c.Count(); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestAdd_And_Count_WithinWindow(t *testing.T) {
	base := time.Now()
	// each call advances by 30 seconds — all within the 5-minute window
	clock := fixedClock(base, 30*time.Second)
	c, _ := newWithClock(Config{Size: 5 * time.Minute}, clock)

	for i := 0; i < 4; i++ {
		c.Add("svc")
	}
	if got := c.Count(); got != 4 {
		t.Errorf("expected 4, got %d", got)
	}
}

func TestAdd_EvictsExpiredEvents(t *testing.T) {
	base := time.Now()
	// each call advances by 2 minutes; window is 3 minutes
	// events at t+0, t+2, t+4 — after t+4 the cutoff is t+1, evicting event at t+0
	clock := fixedClock(base, 2*time.Minute)
	c, _ := newWithClock(Config{Size: 3 * time.Minute}, clock)

	c.Add("svc") // t+0
	c.Add("svc") // t+2
	c.Add("svc") // t+4  — now() is t+6 for Count, cutoff t+3, evicts t+0 and t+2

	if got := c.Count(); got != 1 {
		t.Errorf("expected 1 event within window, got %d", got)
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	c, _ := New(DefaultConfig())
	c.Add("svc")
	c.Add("svc")
	c.Reset()
	if got := c.Count(); got != 0 {
		t.Errorf("expected 0 after reset, got %d", got)
	}
}
