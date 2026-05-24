package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/driftwatch/internal/drift"
)

// WebhookConfig configures an HTTP webhook notifier.
type WebhookConfig struct {
	URL     string
	Level   Level
	Timeout time.Duration
}

type webhookPayload struct {
	Summary  string `json:"summary"`
	HasDrift bool   `json:"has_drift"`
	Drifted  int    `json:"drifted_count"`
}

type webhookNotifier struct {
	cfg    WebhookConfig
	client *http.Client
}

// NewWebhook returns a Notifier that POSTs drift reports to a webhook URL.
func NewWebhook(cfg WebhookConfig) Notifier {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &webhookNotifier{
		cfg:    cfg,
		client: &http.Client{Timeout: timeout},
	}
}

func (w *webhookNotifier) Notify(report drift.Report) error {
	if w.cfg.Level == LevelNone {
		return nil
	}
	if w.cfg.Level == LevelDrift && !report.HasDrift() {
		return nil
	}

	driftedCount := 0
	for _, e := range report.Entries {
		if e.Status == drift.StatusDrifted {
			driftedCount++
		}
	}

	payload := webhookPayload{
		Summary:  report.Summary,
		HasDrift: report.HasDrift(),
		Drifted:  driftedCount,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.cfg.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: webhook post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		// Read a snippet of the response body to include in the error for easier debugging.
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		if len(snippet) > 0 {
			return fmt.Errorf("notify: webhook returned status %d: %s", resp.StatusCode, snippet)
		}
		return fmt.Errorf("notify: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
