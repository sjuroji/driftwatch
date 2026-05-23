package cache

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestCache_SetAndGet_HitBeforeExpiry(t *testing.T) {
	now := time.Now()
	c := New(5 * time.Second)
	c.nowFunc = fixedNow(now)

	c.Set("key", "value")

	got, ok := c.Get("key")
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if got != "value" {
		t.Fatalf("expected 'value', got %v", got)
	}
}

func TestCache_Get_MissAfterExpiry(t *testing.T) {
	now := time.Now()
	c := New(1 * time.Second)
	c.nowFunc = fixedNow(now)
	c.Set("key", "value")

	// Advance time past TTL.
	c.nowFunc = fixedNow(now.Add(2 * time.Second))

	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected cache miss after expiry, got hit")
	}
}

func TestCache_Get_MissingKey(t *testing.T) {
	c := New(time.Minute)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected miss for missing key")
	}
}

func TestCache_Delete_RemovesEntry(t *testing.T) {
	c := New(time.Minute)
	c.Set("key", 42)
	c.Delete("key")

	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected miss after delete")
	}
}

func TestCache_Flush_ClearsAll(t *testing.T) {
	c := New(time.Minute)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Flush()

	if c.Len() != 0 {
		t.Fatalf("expected 0 entries after flush, got %d", c.Len())
	}
}

func TestCache_Len_CountsAll(t *testing.T) {
	c := New(time.Minute)
	c.Set("x", nil)
	c.Set("y", nil)
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestCache_Overwrite_UpdatesEntry(t *testing.T) {
	c := New(time.Minute)
	c.Set("key", "first")
	c.Set("key", "second")

	got, ok := c.Get("key")
	if !ok {
		t.Fatal("expected hit")
	}
	if got != "second" {
		t.Fatalf("expected 'second', got %v", got)
	}
}
