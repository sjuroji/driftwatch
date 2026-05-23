// Package cache provides a lightweight, thread-safe, TTL-based in-memory cache
// used by driftwatch to store live service state snapshots between drift checks.
//
// # Overview
//
// During a scheduled drift-watch run the live fetcher may be called multiple
// times for the same service (e.g. once per manifest). Caching the fetched
// state for a short TTL avoids redundant network calls while still ensuring
// that results are fresh across separate scheduler ticks.
//
// # Usage
//
//	c := cache.New(30 * time.Second)
//	c.Set("auth-service", liveState)
//	if v, ok := c.Get("auth-service"); ok {
//		state := v.(live.State)
//	}
//
// Entries expire automatically based on the TTL supplied to New; no background
// goroutine is required — expiry is checked lazily on Get.
package cache
