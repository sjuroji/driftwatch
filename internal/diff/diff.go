// Package diff provides utilities for computing human-readable field-level
// differences between two drift report snapshots.
package diff

import (
	"fmt"
	"strings"

	"github.com/user/driftwatch/internal/drift"
)

// FieldDiff represents a single field that changed between two snapshots.
type FieldDiff struct {
	Service string
	Field   string
	Old     string
	New     string
}

// String returns a human-readable representation of the diff.
func (f FieldDiff) String() string {
	return fmt.Sprintf("%s.%s: %q -> %q", f.Service, f.Field, f.Old, f.New)
}

// Compare returns the field-level differences between two reports.
// It compares each service entry present in either report.
func Compare(prev, next drift.Report) []FieldDiff {
	var diffs []FieldDiff

	prevIndex := indexByService(prev)
	nextIndex := indexByService(next)

	// Services present in next (added or changed).
	for svc, nextEntry := range nextIndex {
		prevEntry, existed := prevIndex[svc]
		if !existed {
			diffs = append(diffs, FieldDiff{
				Service: svc,
				Field:   "status",
				Old:     "(absent)",
				New:     nextEntry.Status,
			})
			continue
		}
		diffs = append(diffs, compareEntries(svc, prevEntry, nextEntry)...)
	}

	// Services that disappeared.
	for svc, prevEntry := range prevIndex {
		if _, ok := nextIndex[svc]; !ok {
			diffs = append(diffs, FieldDiff{
				Service: svc,
				Field:   "status",
				Old:     prevEntry.Status,
				New:     "(absent)",
			})
		}
	}

	return diffs
}

func compareEntries(svc string, prev, next drift.Entry) []FieldDiff {
	var diffs []FieldDiff
	compare := func(field, a, b string) {
		if a != b {
			diffs = append(diffs, FieldDiff{Service: svc, Field: field, Old: a, New: b})
		}
	}
	compare("status", prev.Status, next.Status)
	compare("image", prev.LiveImage, next.LiveImage)
	compare("replicas",
		fmt.Sprintf("%d", prev.LiveReplicas),
		fmt.Sprintf("%d", next.LiveReplicas))
	return diffs
}

func indexByService(r drift.Report) map[string]drift.Entry {
	m := make(map[string]drift.Entry, len(r.Entries))
	for _, e := range r.Entries {
		m[strings.ToLower(e.Service)] = e
	}
	return m
}
