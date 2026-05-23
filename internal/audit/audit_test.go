package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/internal/audit"
)

func fixedTime() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func newTestLogger(buf *bytes.Buffer) *audit.Logger {
	l := audit.New(buf)
	return l
}

func TestLog_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.Log(audit.LevelInfo, "drift detected", "auth-service", map[string]string{"replicas": "3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var e audit.Event
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if e.Level != audit.LevelInfo {
		t.Errorf("expected level INFO, got %s", e.Level)
	}
	if e.Message != "drift detected" {
		t.Errorf("expected message 'drift detected', got %s", e.Message)
	}
	if e.Service != "auth-service" {
		t.Errorf("expected service 'auth-service', got %s", e.Service)
	}
	if e.Fields["replicas"] != "3" {
		t.Errorf("expected field replicas=3, got %s", e.Fields["replicas"])
	}
}

func TestLog_LevelWarn(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	_ = l.Warn("high replica drift", "payment-service", nil)

	var e audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if e.Level != audit.LevelWarn {
		t.Errorf("expected WARN, got %s", e.Level)
	}
}

func TestLog_LevelError(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	_ = l.Error("fetch failed", "order-service", map[string]string{"reason": "timeout"})

	var e audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if e.Level != audit.LevelError {
		t.Errorf("expected ERROR, got %s", e.Level)
	}
}

func TestNew_NilWriter_UsesStdout(t *testing.T) {
	// Should not panic when nil is passed.
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLog_EmptyService_OmittedFromJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	_ = l.Info("scan started", "", nil)

	raw := strings.TrimSpace(buf.String())
	if strings.Contains(raw, `"service"`) {
		t.Errorf("expected service field to be omitted, got: %s", raw)
	}
}
