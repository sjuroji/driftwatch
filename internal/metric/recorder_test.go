package metric_test

import (
	"testing"
	"time"

	"driftwatch/internal/drift"
	"driftwatch/internal/metric"
)

func makeReport(statuses []drift.Status) drift.Report {
	entries := make([]drift.Entry, len(statuses))
	for i, s := range statuses {
		entries[i] = drift.Entry{
			ServiceName: "svc",
			Status:      s,
		}
	}
	return drift.Report{Entries: entries}
}

func TestFromReport_AllInSync(t *testing.T) {
	report := makeReport([]drift.Status{
		drift.StatusInSync,
		drift.StatusInSync,
	})
	r := metric.FromReport(report, 50*time.Millisecond)

	if r.TotalServices != 2 {
		t.Errorf("expected TotalServices=2, got %d", r.TotalServices)
	}
	if r.InSync != 2 {
		t.Errorf("expected InSync=2, got %d", r.InSync)
	}
	if r.Drifted != 0 {
		t.Errorf("expected Drifted=0, got %d", r.Drifted)
	}
	if r.DurationMs != 50 {
		t.Errorf("expected DurationMs=50, got %d", r.DurationMs)
	}
}

func TestFromReport_MixedStatuses(t *testing.T) {
	report := makeReport([]drift.Status{
		drift.StatusInSync,
		drift.StatusDrifted,
		drift.StatusDrifted,
	})
	r := metric.FromReport(report, 100*time.Millisecond)

	if r.Drifted != 2 {
		t.Errorf("expected Drifted=2, got %d", r.Drifted)
	}
	if r.InSync != 1 {
		t.Errorf("expected InSync=1, got %d", r.InSync)
	}
}

func TestFromReport_EmptyReport(t *testing.T) {
	r := metric.FromReport(drift.Report{}, 10*time.Millisecond)
	if r.TotalServices != 0 || r.InSync != 0 || r.Drifted != 0 {
		t.Errorf("unexpected values for empty report: %+v", r)
	}
}
