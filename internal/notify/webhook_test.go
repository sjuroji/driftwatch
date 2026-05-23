package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/driftwatch/internal/notify"
)

func TestWebhookNotifier_PostsOnDrift(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewWebhook(notify.WebhookConfig{
		URL:   server.URL,
		Level: notify.LevelDrift,
	})
	if err := n.Notify(driftReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["has_drift"] != true {
		t.Errorf("expected has_drift=true, got %v", received["has_drift"])
	}
}

func TestWebhookNotifier_SilentOnSyncWithDriftLevel(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewWebhook(notify.WebhookConfig{
		URL:   server.URL,
		Level: notify.LevelDrift,
	})
	if err := n.Notify(syncReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected webhook not to be called for in-sync report")
	}
}

func TestWebhookNotifier_ErrorOnBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := notify.NewWebhook(notify.WebhookConfig{
		URL:     server.URL,
		Level:   notify.LevelAll,
		Timeout: 2 * time.Second,
	})
	if err := n.Notify(syncReport()); err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestWebhookNotifier_ErrorOnBadURL(t *testing.T) {
	n := notify.NewWebhook(notify.WebhookConfig{
		URL:     "http://127.0.0.1:0/no-server",
		Level:   notify.LevelAll,
		Timeout: 500 * time.Millisecond,
	})
	if err := n.Notify(driftReport()); err == nil {
		t.Error("expected connection error")
	}
}
