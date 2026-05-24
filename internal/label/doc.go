// Package label provides label selector parsing, matching, filtering, and
// grouping utilities for driftwatch service manifests.
//
// A Selector is a set of key=value requirements. A manifest or live-state
// entry matches a selector when it carries all required labels with the
// expected values. Key comparison is case-insensitive.
//
// Typical usage:
//
//	sel, err := label.ParseSelector("env=prod,team=platform")
//	if err != nil { ... }
//
//	filtered := label.Filter(manifests, sel)
//	groups   := label.GroupBy(manifests, "team")
package label
