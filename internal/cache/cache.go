// Package cache provides a simple in-memory cache for live state snapshots,
// reducing redundant fetches during a single drift-watch run.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached value along with its expiry time.
type Entry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// Cache is a thread-safe, TTL-based in-memory store.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]Entry
	ttl     time.Duration
	nowFunc func() time.Time // injectable for testing
}

// New creates a Cache with the given TTL for all entries.
func New(ttl time.Duration) *Cache {
	return &Cache{
		items:   make(map[string]Entry),
		ttl:     ttl,
		nowFunc: time.Now,
	}
}

// Set stores value under key, overwriting any existing entry.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Entry{
		Value:     value,
		ExpiresAt: c.nowFunc().Add(c.ttl),
	}
}

// Get retrieves the value for key. Returns (value, true) if the entry exists
// and has not expired, otherwise (nil, false).
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if c.nowFunc().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Value, true
}

// Delete removes key from the cache if present.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]Entry)
}

// Len returns the number of entries currently in the cache (including expired).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
