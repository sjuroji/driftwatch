package trend_test

import (
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/trend"
)

func makeReport(statuses map[string]drift.Status) drift.Report {
	entries := make([]drift.Entry, 0, len(statuses))
	for svc, status := range statuses {
		entries = append(entries, drift.Entry{Service: svc, Status: status})
	}
	return drift.Report{Entries: entries}
}

func TestAnalyse_NoReports_ReturnsError(t *testing.T) {
	_, err := trend.Analyse(nil, time.Hour)
	if err == nil {
		t.Fatal("expected error for empty reports, got nil")
	}
}

func TestAnalyse_AllInSync_ZeroDriftRate(t *testing.T) {
	reports := []drift.Report{
		makeReport(map[string]drift.Status{"api": drift.StatusInSync}),
		makeReport(map[string]drift.Status{"api": drift.StatusInSync}),
	}

	r, err := trend.Analyse(reports, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(r.Services))
	}
	if r.Services[0].DriftRate != 0 {
		t.Errorf("expected drift rate 0, got %f", r.Services[0].DriftRate)
	}
	if r.Services[0].Direction != trend.DirectionStable {
		t.Errorf("expected stable direction, got %s", r.Services[0].Direction)
	}
}

func TestAnalyse_AlwaysDrifted_FullRate(t *testing.T) {
	reports := []drift.Report{
		makeReport(map[string]drift.Status{"worker": drift.StatusDrifted}),
		makeReport(map[string]drift.Status{"worker": drift.StatusDrifted}),
		makeReport(map[string]drift.Status{"worker": drift.StatusDrifted}),
		makeReport(map[string]drift.Status{"worker": drift.StatusDrifted}),
	}

	r, err := trend.Analyse(reports, 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Services[0].DriftRate != 1.0 {
		t.Errorf("expected drift rate 1.0, got %f", r.Services[0].DriftRate)
	}
	if r.Services[0].Direction != trend.DirectionStable {
		t.Errorf("expected stable (consistently drifted), got %s", r.Services[0].Direction)
	}
}

func TestAnalyse_Worsening(t *testing.T) {
	// early half: in-sync; late half: drifted → worsening
	reports := []drift.Report{
		makeReport(map[string]drift.Status{"svc": drift.StatusInSync}),
		makeReport(map[string]drift.Status{"svc": drift.StatusInSync}),
		makeReport(map[string]drift.Status{"svc": drift.StatusDrifted}),
		makeReport(map[string]drift.Status{"svc": drift.StatusDrifted}),
	}

	r, err := trend.Analyse(reports, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Services[0].Direction != trend.DirectionWorsening {
		t.Errorf("expected worsening, got %s", r.Services[0].Direction)
	}
}

func TestAnalyse_Improving(t *testing.T) {
	// early half: drifted; late half: in-sync → improving
	reports := []drift.Report{
		makeReport(map[string]drift.Status{"svc": drift.StatusDrifted}),
		makeReport(map[string]drift.Status{"svc": drift.StatusDrifted}),
		makeReport(map[string]drift.Status{"svc": drift.StatusInSync}),
		makeReport(map[string]drift.Status{"svc": drift.StatusInSync}),
	}

	r, err := trend.Analyse(reports, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Services[0].Direction != trend.DirectionImproving {
		t.Errorf("expected improving, got %s", r.Services[0].Direction)
	}
}

func TestAnalyse_MultipleServices_SortedByDriftRate(t *testing.T) {
	reports := []drift.Report{
		makeReport(map[string]drift.Status{
			"low":  drift.StatusInSync,
			"high": drift.StatusDrifted,
		}),
	}

	r, err := trend.Analyse(reports, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(r.Services))
	}
	if r.Services[0].Service != "high" {
		t.Errorf("expected highest drift rate first, got %s", r.Services[0].Service)
	}
}
