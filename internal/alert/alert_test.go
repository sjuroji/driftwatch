package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/alert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := alert.DefaultConfig()
	if cfg.WarnThreshold != 1 {
		t.Errorf("expected WarnThreshold=1, got %d", cfg.WarnThreshold)
	}
	if cfg.CritThreshold != 5 {
		t.Errorf("expected CritThreshold=5, got %d", cfg.CritThreshold)
	}
}

func TestEvaluate_BelowWarn_ReturnsNil(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(alert.Config{WarnThreshold: 2, CritThreshold: 5}, &buf)
	result := a.Evaluate(1)
	if result != nil {
		t.Errorf("expected nil alert, got %+v", result)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestEvaluate_AtWarnThreshold(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(alert.Config{WarnThreshold: 2, CritThreshold: 5}, &buf)
	result := a.Evaluate(2)
	if result == nil {
		t.Fatal("expected alert, got nil")
	}
	if result.Level != alert.LevelWarn {
		t.Errorf("expected level %q, got %q", alert.LevelWarn, result.Level)
	}
	if result.Drifted != 2 {
		t.Errorf("expected drifted=2, got %d", result.Drifted)
	}
	if !strings.Contains(buf.String(), "warn") {
		t.Errorf("expected output to contain 'warn', got: %s", buf.String())
	}
}

func TestEvaluate_AtCritThreshold(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(alert.Config{WarnThreshold: 2, CritThreshold: 5}, &buf)
	result := a.Evaluate(5)
	if result == nil {
		t.Fatal("expected alert, got nil")
	}
	if result.Level != alert.LevelCrit {
		t.Errorf("expected level %q, got %q", alert.LevelCrit, result.Level)
	}
	if !strings.Contains(buf.String(), "critical") {
		t.Errorf("expected output to contain 'critical', got: %s", buf.String())
	}
}

func TestEvaluate_AboveCrit(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(alert.Config{WarnThreshold: 1, CritThreshold: 3}, &buf)
	result := a.Evaluate(10)
	if result == nil {
		t.Fatal("expected alert, got nil")
	}
	if result.Level != alert.LevelCrit {
		t.Errorf("expected critical, got %q", result.Level)
	}
}

func TestNew_NilWriter_UsesStderr(t *testing.T) {
	// Should not panic when out is nil.
	a := alert.New(alert.DefaultConfig(), nil)
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}
