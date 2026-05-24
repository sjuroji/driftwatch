package severity_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/severity"
)

func makeReport(driftedCount, syncedCount int) drift.Report {
	entries := make([]drift.Entry, 0, driftedCount+syncedCount)
	for i := 0; i < driftedCount; i++ {
		entries = append(entries, drift.Entry{
			Service: fmt.Sprintf("svc-drift-%d", i),
			Status:  drift.StatusDrifted,
		})
	}
	for i := 0; i < syncedCount; i++ {
		entries = append(entries, drift.Entry{
			Service: fmt.Sprintf("svc-sync-%d", i),
			Status:  drift.StatusInSync,
		})
	}
	return drift.Report{Entries: entries}
}

func TestDefaultConfig(t *testing.T) {
	cfg := severity.DefaultConfig()
	if cfg.LowAt != 1 || cfg.MediumAt != 3 || cfg.HighAt != 6 || cfg.CriticalAt != 10 {
		t.Fatalf("unexpected default config: %+v", cfg)
	}
}

func TestEvaluate_NoDrift_ReturnsNone(t *testing.T) {
	e := severity.New(severity.DefaultConfig())
	level := e.Evaluate(makeReport(0, 5))
	if level != severity.LevelNone {
		t.Fatalf("expected none, got %s", level)
	}
}

func TestEvaluate_OneDrifted_ReturnsLow(t *testing.T) {
	e := severity.New(severity.DefaultConfig())
	level := e.Evaluate(makeReport(1, 4))
	if level != severity.LevelLow {
		t.Fatalf("expected low, got %s", level)
	}
}

func TestEvaluate_AtMediumThreshold(t *testing.T) {
	e := severity.New(severity.DefaultConfig())
	level := e.Evaluate(makeReport(3, 2))
	if level != severity.LevelMedium {
		t.Fatalf("expected medium, got %s", level)
	}
}

func TestEvaluate_AtHighThreshold(t *testing.T) {
	e := severity.New(severity.DefaultConfig())
	level := e.Evaluate(makeReport(6, 0))
	if level != severity.LevelHigh {
		t.Fatalf("expected high, got %s", level)
	}
}

func TestEvaluate_AtCriticalThreshold(t *testing.T) {
	e := severity.New(severity.DefaultConfig())
	level := e.Evaluate(makeReport(10, 0))
	if level != severity.LevelCritical {
		t.Fatalf("expected critical, got %s", level)
	}
}

func TestEvaluate_ZeroValueConfig_UsesDefaults(t *testing.T) {
	e := severity.New(severity.Config{})
	level := e.Evaluate(makeReport(0, 3))
	if level != severity.LevelNone {
		t.Fatalf("expected none, got %s", level)
	}
}

func TestEvaluate_CustomConfig(t *testing.T) {
	cfg := severity.Config{LowAt: 2, MediumAt: 5, HighAt: 8, CriticalAt: 15}
	e := severity.New(cfg)

	if got := e.Evaluate(makeReport(1, 0)); got != severity.LevelNone {
		t.Fatalf("expected none, got %s", got)
	}
	if got := e.Evaluate(makeReport(2, 0)); got != severity.LevelLow {
		t.Fatalf("expected low, got %s", got)
	}
	if got := e.Evaluate(makeReport(15, 0)); got != severity.LevelCritical {
		t.Fatalf("expected critical, got %s", got)
	}
}
