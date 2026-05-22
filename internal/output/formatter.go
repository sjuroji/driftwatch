package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Format controls the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Formatter writes a drift report to a writer in a specified format.
type Formatter interface {
	Write(w io.Writer, report drift.Report) error
}

// New returns a Formatter for the given format.
func New(f Format) (Formatter, error) {
	switch f {
	case FormatText:
		return &TextFormatter{}, nil
	case FormatJSON:
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %q", f)
	}
}

// TextFormatter renders drift reports as human-readable text.
type TextFormatter struct{}

func (t *TextFormatter) Write(w io.Writer, report drift.Report) error {
	fmt.Fprintf(w, "Drift Report\n")
	fmt.Fprintf(w, "%s\n", strings.Repeat("=", 40))
	fmt.Fprintf(w, "Services checked : %d\n", report.Summary.TotalServices)
	fmt.Fprintf(w, "In sync          : %d\n", report.Summary.InSync)
	fmt.Fprintf(w, "Drifted          : %d\n", report.Summary.Drifted)
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 40))

	for _, entry := range report.Entries {
		if entry.InSync {
			fmt.Fprintf(w, "[OK]   %s\n", entry.ServiceName)
			continue
		}
		fmt.Fprintf(w, "[DRIFT] %s\n", entry.ServiceName)
		for _, d := range entry.Diffs {
			fmt.Fprintf(w, "        field=%s want=%v got=%v\n", d.Field, d.Want, d.Got)
		}
	}
	return nil
}
