// Package digest provides lightweight content-addressable hashing for
// driftwatch service states.
//
// Instead of performing a full structural diff on every poll cycle,
// callers can compute a SHA-256 digest of the current manifest or live
// state and compare it against a previously stored digest.  Only when
// the digests differ is a more expensive field-level comparison needed.
//
// Basic usage:
//
//	prev, _ := digest.Compute("auth-service", previousState)
//	curr, _ := digest.Compute("auth-service", currentState)
//	if digest.Changed(prev, curr) {
//		// run full drift detection
//	}
//
// ComputeAll and Index are convenience helpers for working with multiple
// services at once.
package digest
