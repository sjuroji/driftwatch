// Package ratelimit provides a thread-safe, per-key token-bucket rate
// limiter intended to guard high-frequency paths in driftwatch — such
// as notification dispatch and webhook delivery — from overwhelming
// downstream systems during sustained drift events.
//
// # Usage
//
//	cfg := ratelimit.DefaultConfig() // 5 events / minute per key
//	limiter, err := ratelimit.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if limiter.Allow(serviceName) {
//		// safe to notify
//	}
//
// Each unique key (typically a service name) maintains its own
// independent counter that resets at the end of each Interval window.
package ratelimit
