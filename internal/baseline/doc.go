// Package baseline provides storage for known-good service configuration
// states used as reference points during drift detection.
//
// A baseline captures the live state of a service at a point in time when
// the operator has declared it to be correct. Subsequent drift checks can
// compare the current live state against the stored baseline rather than
// (or in addition to) the declared manifest.
//
// Usage:
//
//	store, err := baseline.NewStore("/var/lib/driftwatch/baselines")
//	if err != nil { ... }
//
//	// Capture the current state as the new baseline.
//	if err := store.Save("auth-service", liveState, time.Now()); err != nil { ... }
//
//	// Retrieve a previously saved baseline.
//	entry, err := store.Load("auth-service")
//	if errors.Is(err, os.ErrNotExist) { /* no baseline yet */ }
package baseline
