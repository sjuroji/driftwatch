package diff_test

import (
	"testing"

	"github.com/user/driftwatch/internal/diff"
	"github.com/user/driftwatch/internal/drift"
)

func baseReport() drift.Report {
	return drift.Report{
		Entries: []drift.Entry{
			{
				Service:      "auth-service",
				Status:       "in-sync",
				LiveImage:    "auth:v1",
				LiveReplicas: 2,
			},
			{
				Service:      "api-gateway",
				Status:       "drift",
				LiveImage:    "api:v2",
				LiveReplicas: 3,
			},
		},
	}
}

func TestCompare_NoDiff(t *testing.T) {
	r := baseReport()
	diffs := diff.Compare(r, r)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs, got %d: %v", len(diffs), diffs)
	}
}

func TestCompare_StatusChange(t *testing.T) {
	prev := baseReport()
	next := baseReport()
	next.Entries[0].Status = "drift"

	diffs := diff.Compare(prev, next)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Field != "status" || diffs[0].Old != "in-sync" || diffs[0].New != "drift" {
		t.Errorf("unexpected diff: %v", diffs[0])
	}
}

func TestCompare_ImageChange(t *testing.T) {
	prev := baseReport()
	next := baseReport()
	next.Entries[1].LiveImage = "api:v3"

	diffs := diff.Compare(prev, next)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Field != "image" {
		t.Errorf("expected image diff, got field=%q", diffs[0].Field)
	}
}

func TestCompare_ServiceAdded(t *testing.T) {
	prev := baseReport()
	next := baseReport()
	next.Entries = append(next.Entries, drift.Entry{
		Service: "new-svc",
		Status:  "in-sync",
	})

	diffs := diff.Compare(prev, next)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Old != "(absent)" {
		t.Errorf("expected absent old value, got %q", diffs[0].Old)
	}
}

func TestCompare_ServiceRemoved(t *testing.T) {
	prev := baseReport()
	next := baseReport()
	next.Entries = next.Entries[:1]

	diffs := diff.Compare(prev, next)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].New != "(absent)" {
		t.Errorf("expected absent new value, got %q", diffs[0].New)
	}
}

func TestFieldDiff_String(t *testing.T) {
	d := diff.FieldDiff{Service: "svc", Field: "image", Old: "v1", New: "v2"}
	got := d.String()
	expected := `svc.image: "v1" -> "v2"`
	if got != expected {
		t.Errorf("String() = %q, want %q", got, expected)
	}
}
