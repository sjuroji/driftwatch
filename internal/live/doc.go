// Package live provides types and utilities for fetching the current
// runtime state of deployed services.
//
// The core abstraction is the [Fetcher] interface, which allows different
// backends (Kubernetes, Docker, static stubs) to supply live service state.
// [LiveState] mirrors the fields declared in a [manifest.Manifest] so that
// the drift detector can compare declared vs. actual configuration.
//
// For testing and local development, [StaticFetcher] provides an in-memory
// implementation backed by a predefined map of states.
//
// Typical usage:
//
//	fetcher := &live.StaticFetcher{States: myStates}
//	states, errs := live.FetchAll(fetcher, manifests)
package live
