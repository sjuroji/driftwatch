package history_test

import (
	"os"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/history"
)

func sampleReport(drifted bool) drift.Report {
	status := drift.StatusInSync
	if drifted {
		status = drift.StatusDrifted
	}
	return drift.Report{
		Results: []drift.Result{
			{
				Service: "auth-service",
				Status:  status,
				Diffs:   nil,
			},
		},
		Summary: drift.Summary{Total: 1, Drifted: func() int {
			if drifted {
				return 1
			}
			return 0
		}()},
	}
}

func TestNewStore_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	sub := dir + "/nested/store"

	_, err := history.NewStore(sub)
	if err != nil {
		t.Fatalf("NewStore: unexpected error: %v", err)
	}

	if _, err := os.Stat(sub); os.IsNotExist(err) {
		t.Errorf("expected directory %s to be created", sub)
	}
}

func TestStore_SaveAndList_RoundTrip(t *testing.T) {
	store, err := history.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	report := sampleReport(false)
	if err := store.Save(report); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Report.Summary.Total != 1 {
		t.Errorf("expected Summary.Total=1, got %d", entries[0].Report.Summary.Total)
	}
}

func TestStore_List_MultipleEntries_OldestFirst(t *testing.T) {
	store, err := history.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := store.Save(sampleReport(i%2 == 0)); err != nil {
			t.Fatalf("Save[%d]: %v", i, err)
		}
	}

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].Timestamp.Before(entries[i-1].Timestamp) {
			t.Errorf("entries not sorted oldest-first at index %d", i)
		}
	}
}

func TestStore_List_EmptyStore(t *testing.T) {
	store, err := history.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
