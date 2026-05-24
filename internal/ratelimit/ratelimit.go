// Package ratelimit provides a simple token-bucket rate limiter
// for controlling how frequently drift checks or notifications are
// emitted per service.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Config holds rate limiter settings.
type Config struct {
	// Rate is the number of allowed events per Interval.
	Rate int
	// Interval is the window over which Rate events are allowed.
	Interval time.Duration
}

// DefaultConfig returns a sensible default: 5 events per minute.
func DefaultConfig() Config {
	return Config{
		Rate:     5,
		Interval: time.Minute,
	}
}

// bucket tracks usage for a single key.
type bucket struct {
	count     int
	windowEnd time.Time
}

// Limiter enforces per-key rate limits.
type Limiter struct {
	cfg     Config
	mu      sync.Mutex
	buckets map[string]*bucket
	now     func() time.Time
}

// New creates a Limiter with the given Config.
func New(cfg Config) (*Limiter, error) {
	if cfg.Rate <= 0 {
		return nil, fmt.Errorf("ratelimit: Rate must be > 0, got %d", cfg.Rate)
	}
	if cfg.Interval <= 0 {
		return nil, fmt.Errorf("ratelimit: Interval must be > 0, got %v", cfg.Interval)
	}
	return &Limiter{
		cfg:     cfg,
		buckets: make(map[string]*bucket),
		now:     time.Now,
	}, nil
}

// Allow reports whether the event for key is within the rate limit.
// It increments the counter for the current window.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	b, ok := l.buckets[key]
	if !ok || now.After(b.windowEnd) {
		l.buckets[key] = &bucket{
			count:     1,
			windowEnd: now.Add(l.cfg.Interval),
		}
		return true
	}

	if b.count >= l.cfg.Rate {
		return false
	}
	b.count++
	return true
}

// Reset clears the rate limit state for the given key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}
