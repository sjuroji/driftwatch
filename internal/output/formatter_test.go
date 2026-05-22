package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/output"
)

func sampleReport() drift.Report {
	return drift.Report{
		Summary: drift.Summary{
			TotalServices: 2,
			InSync:        1,
			Drifted:       1,
		},
		Entries: []drift.Entry{
			{ServiceName: "auth-service", InSync: true},
			{
				ServiceName: "payment-service",
				InSync:      false,
				Diffs: []drift.Diff{
					{Field: "replicas", Want: 3, Got: 1},
				},
			},
		},
	}
}

func TestNew_UnknownFormat(t *testing.T) {
	_, err := output.New("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestTextFormatter_InSyncLine(t *testing.T) {
	f, _ := output.New(output.FormatText)
	var buf bytes.Buffer
	if err := f.Write(&buf, sampleReport()); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if !strings.Contains(buf.String(), "[OK]   auth-service") {
		t.Errorf("expected in-sync line for auth-service, got:\n%s", buf.String())
	}
}

func TestTextFormatter_DriftLine(t *testing.T) {
	f, _ := output.New(output.FormatText)
	var buf bytes.Buffer
	f.Write(&buf, sampleReport())
	out := buf.String()
	if !strings.Contains(out, "[DRIFT] payment-service") {
		t.Errorf("expected drift line for payment-service, got:\n%s", out)
	}
	if !strings.Contains(out, "field=replicas") {
		t.Errorf("expected replicas diff in output, got:\n%s", out)
	}
}

func TestJSONFormatter_ValidJSON(t *testing.T) {
	f, _ := output.New(output.FormatJSON)
	var buf bytes.Buffer
	if err := f.Write(&buf, sampleReport()); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	var got drift.Report
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, buf.String())
	}
	if got.Summary.Drifted != 1 {
		t.Errorf("expected Drifted=1, got %d", got.Summary.Drifted)
	}
}

func TestTextFormatter_SummaryCounts(t *testing.T) {
	f, _ := output.New(output.FormatText)
	var buf bytes.Buffer
	f.Write(&buf, sampleReport())
	out := buf.String()
	if !strings.Contains(out, "Services checked : 2") {
		t.Errorf("expected total count in summary, got:\n%s", out)
	}
}
