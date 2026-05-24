// Package window provides a sliding time-window counter used to track
// how many drift events have occurred within a configurable duration.
package window

import (
	"fmt"
	"sync"
	"time"
)

// Config holds the configuration for a sliding window.
type Config struct {
	// Size is the duration of the window (e.g. 5 minutes).
	Size time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Size: 5 * time.Minute,
	}
}

// Counter tracks timestamped events within a sliding time window.
type Counter struct {
	mu     sync.Mutex
	cfg    Config
	events []time.Time
	now    func() time.Time
}

// New creates a Counter using the provided Config.
// Returns an error if the window size is not positive.
func New(cfg Config) (*Counter, error) {
	return newWithClock(cfg, time.Now)
}

func newWithClock(cfg Config, now func() time.Time) (*Counter, error) {
	if cfg.Size <= 0 {
		return nil, fmt.Errorf("window: size must be positive, got %s", cfg.Size)
	}
	return &Counter{cfg: cfg, now: now}, nil
}

// Add records a new event at the current time.
func (c *Counter) Add(service string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, c.now())
	c.evict()
}

// Count returns the number of events recorded within the current window.
func (c *Counter) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict()
	return len(c.events)
}

// Reset clears all recorded events.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = c.events[:0]
}

// evict removes events older than the window size.
// Must be called with c.mu held.
func (c *Counter) evict() {
	cutoff := c.now().Add(-c.cfg.Size)
	i := 0
	for i < len(c.events) && c.events[i].Before(cutoff) {
		i++
	}
	c.events = c.events[i:]
}
