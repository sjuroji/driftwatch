package rollup_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/rollup"
)

func makeReport(statuses map[string]drift.Status) drift.Report {
	entries := make([]drift.Entry, 0, len(statuses))
	for svc, status := range statuses {
		entries = append(entries, drift.Entry{Service: svc, Status: status})
	}
	return drift.Report{Entries: entries}
}

func TestAggregate_NoReports_ReturnsError(t *testing.T) {
	_, err := rollup.Aggregate(nil)
	if err == nil {
		t.Fatal("expected error for empty reports, got nil")
	}
}

func TestAggregate_AllInSync_ZeroDriftRate(t *testing.T) {
	reports := []drift.Report{
		makeReport(map[string]drift.Status{"api": drift.StatusInSync, "worker": drift.StatusInSync}),
		makeReport(map[string]drift.Status{"api": drift.StatusInSync, "worker": drift.StatusInSync}),
	}
	res, err := rollup.Aggregate(reports)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TotalReports != 2 {
		t.Errorf("TotalReports = %d, want 2", res.TotalReports)
	}
	for _, s := range res.Services {
		if s.DriftRate != 0.0 {
			t.Errorf("service %s: DriftRate = %.2f, want 0.0", s.Service, s.DriftRate)
		}
	}
}

func TestAggregate_AlwaysDrifted_FullDriftRate(t *testing.T) {
	reports := []drift.Report{
		makeReport(map[string]drift.Status{"auth": drift.StatusDrift}),
		makeReport(map[string]drift.Status{"auth": drift.StatusDrift}),
		makeReport(map[string]drift.Status{"auth": drift.StatusDrift}),
	}
	res, err := rollup.Aggregate(reports)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(res.Services))
	}
	s := res.Services[0]
	if s.DriftRate != 1.0 {
		t.Errorf("DriftRate = %.2f, want 1.0", s.DriftRate)
	}
	if s.DriftCount != 3 {
		t.Errorf("DriftCount = %d, want 3", s.DriftCount)
	}
}

func TestAggregate_MixedDrift_CorrectRate(t *testing.T) {
	reports := []drift.Report{
		makeReport(map[string]drift.Status{"svc": drift.StatusDrift}),
		makeReport(map[string]drift.Status{"svc": drift.StatusInSync}),
		makeReport(map[string]drift.Status{"svc": drift.StatusDrift}),
		makeReport(map[string]drift.Status{"svc": drift.StatusInSync}),
	}
	res, err := rollup.Aggregate(reports)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := res.Services[0]
	const want = 0.5
	if s.DriftRate != want {
		t.Errorf("DriftRate = %.2f, want %.2f", s.DriftRate, want)
	}
}

func TestAggregate_SortedByDriftRateDescending(t *testing.T) {
	reports := []drift.Report{
		makeReport(map[string]drift.Status{"low": drift.StatusInSync, "high": drift.StatusDrift}),
		makeReport(map[string]drift.Status{"low": drift.StatusDrift, "high": drift.StatusDrift}),
		makeReport(map[string]drift.Status{"low": drift.StatusInSync, "high": drift.StatusDrift}),
	}
	res, err := rollup.Aggregate(reports)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Services[0].Service != "high" {
		t.Errorf("first service = %q, want \"high\"", res.Services[0].Service)
	}
}
