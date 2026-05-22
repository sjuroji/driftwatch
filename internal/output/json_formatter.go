package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/driftwatch/internal/drift"
)

// JSONFormatter renders drift reports as JSON.
type JSONFormatter struct{}

func (j *JSONFormatter) Write(w io.Writer, report drift.Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("json formatter: %w", err)
	}
	return nil
}
