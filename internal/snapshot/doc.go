// Package snapshot provides point-in-time capture and retrieval of live
// service state for the driftwatch tool.
//
// A [Snapshot] holds a timestamped slice of [Entry] values, each describing
// the observed image, replica count, and labels for a single service.
//
// Snapshots are persisted to disk by [Store] as JSON files named by their
// capture timestamp (e.g. 20240601T120000Z.json). [Store.Latest] returns the
// most recently saved snapshot, which can be compared against a freshly
// captured state to detect changes between runs.
//
// Typical usage:
//
//	store, err := snapshot.NewStore("/var/lib/driftwatch/snapshots")
//	if err != nil { ... }
//
//	snap := snapshot.Snapshot{
//	    CapturedAt: time.Now(),
//	    Entries:    entries,
//	}
//	if err := store.Save(snap); err != nil { ... }
package snapshot
