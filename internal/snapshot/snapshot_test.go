package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/snapshot"
)

func fixedTime() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func sampleSnapshot(t time.Time) snapshot.Snapshot {
	return snapshot.Snapshot{
		CapturedAt: t,
		Entries: []snapshot.Entry{
			{Service: "auth", Image: "auth:v1", Replicas: 2, Labels: map[string]string{"env": "prod"}},
			{Service: "api", Image: "api:v3", Replicas: 4},
		},
	}
}

func TestNewStore_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	subDir := dir + "/snapshots"

	_, err := snapshot.NewStore(subDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestSave_And_Latest_RoundTrip(t *testing.T) {
	store, err := snapshot.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	snap := sampleSnapshot(fixedTime())
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}

	if got.CapturedAt.UTC() != snap.CapturedAt.UTC() {
		t.Errorf("CapturedAt: got %v, want %v", got.CapturedAt, snap.CapturedAt)
	}
	if len(got.Entries) != len(snap.Entries) {
		t.Fatalf("entries len: got %d, want %d", len(got.Entries), len(snap.Entries))
	}
	if got.Entries[0].Service != "auth" {
		t.Errorf("first entry service: got %q, want %q", got.Entries[0].Service, "auth")
	}
}

func TestLatest_ReturnsNewest(t *testing.T) {
	store, err := snapshot.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	old := sampleSnapshot(fixedTime())
	newer := sampleSnapshot(fixedTime().Add(time.Hour))
	newer.Entries[0].Image = "auth:v2"

	for _, s := range []snapshot.Snapshot{old, newer} {
		if err := store.Save(s); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}

	got, err := store.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if got.Entries[0].Image != "auth:v2" {
		t.Errorf("expected newest snapshot; image: got %q, want %q", got.Entries[0].Image, "auth:v2")
	}
}

func TestLatest_EmptyStore_ReturnsError(t *testing.T) {
	store, err := snapshot.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	_, err = store.Latest()
	if err == nil {
		t.Fatal("expected error for empty store, got nil")
	}
}
