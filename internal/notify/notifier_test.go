package notify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/notify"
)

func driftReport() drift.Report {
	return drift.Report{
		Summary: "1 service(s) drifted",
		Entries: []drift.Entry{
			{ServiceName: "auth", Status: drift.StatusDrifted},
		},
	}
}

func syncReport() drift.Report {
	return drift.Report{
		Summary: "all services in sync",
		Entries: []drift.Entry{
			{ServiceName: "auth", Status: drift.StatusInSync},
		},
	}
}

func TestNotify_LevelNone_Silent(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.Config{Level: notify.LevelNone, Writer: &buf})
	if err := n.Notify(driftReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestNotify_LevelDrift_WritesOnDrift(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.Config{Level: notify.LevelDrift, Writer: &buf})
	if err := n.Notify(driftReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "drifted") {
		t.Errorf("expected drift message, got %q", buf.String())
	}
}

func TestNotify_LevelDrift_SilentOnSync(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.Config{Level: notify.LevelDrift, Writer: &buf})
	if err := n.Notify(syncReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output on sync, got %q", buf.String())
	}
}

func TestNotify_LevelAll_WritesOnSync(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.Config{Level: notify.LevelAll, Writer: &buf})
	if err := n.Notify(syncReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[driftwatch]") {
		t.Errorf("expected output, got %q", buf.String())
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := notify.DefaultConfig()
	if cfg.Level != notify.LevelDrift {
		t.Errorf("expected default level %q, got %q", notify.LevelDrift, cfg.Level)
	}
	if cfg.Writer == nil {
		t.Error("expected non-nil writer")
	}
}
