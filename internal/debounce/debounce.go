// Package debounce provides a mechanism to suppress repeated drift
// notifications within a configurable cooldown window. This prevents
// alert fatigue when a service repeatedly drifts between checks.
package debounce

import (
	"sync"
	"time"
)

// DefaultCooldown is the default minimum duration between notifications
// for the same service.
const DefaultCooldown = 5 * time.Minute

// Clock is a function that returns the current time, allowing tests to
// inject a fixed clock.
type Clock func() time.Time

// Debouncer tracks the last notification time per service and suppresses
// repeated firings within the cooldown window.
type Debouncer struct {
	mu       sync.Mutex
	last     map[string]time.Time
	cooldown time.Duration
	now      Clock
}

// New returns a Debouncer with the given cooldown. If cooldown is zero
// DefaultCooldown is used. If clock is nil, time.Now is used.
func New(cooldown time.Duration, clock Clock) *Debouncer {
	if cooldown <= 0 {
		cooldown = DefaultCooldown
	}
	if clock == nil {
		clock = time.Now
	}
	return &Debouncer{
		last:     make(map[string]time.Time),
		cooldown: cooldown,
		now:      clock,
	}
}

// Allow returns true if the service should fire a notification now,
// i.e. it has never fired or its cooldown has elapsed. Calling Allow
// with a true result records the current time as the last fire time.
func (d *Debouncer) Allow(service string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if last, ok := d.last[service]; ok {
		if now.Sub(last) < d.cooldown {
			return false
		}
	}
	d.last[service] = now
	return true
}

// Reset clears the recorded fire time for a service, allowing it to
// fire immediately on the next call to Allow.
func (d *Debouncer) Reset(service string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.last, service)
}

// ResetAll clears all recorded fire times.
func (d *Debouncer) ResetAll() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.last = make(map[string]time.Time)
}
