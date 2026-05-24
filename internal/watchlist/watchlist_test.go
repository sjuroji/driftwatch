package watchlist_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/watchlist"
)

func makeEntry(name string, tags ...string) watchlist.Entry {
	return watchlist.Entry{Name: name, Tags: tags}
}

func TestAdd_And_Contains(t *testing.T) {
	w := watchlist.New()
	if err := w.Add(makeEntry("auth-service", "prod")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !w.Contains("auth-service") {
		t.Error("expected auth-service to be present")
	}
}

func TestAdd_CaseInsensitive(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(makeEntry("Auth-Service"))
	if !w.Contains("auth-service") {
		t.Error("lookup should be case-insensitive")
	}
}

func TestAdd_EmptyName_ReturnsError(t *testing.T) {
	w := watchlist.New()
	if err := w.Add(watchlist.Entry{Name: "   "}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestRemove_ExistingEntry(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(makeEntry("billing"))
	if !w.Remove("billing") {
		t.Error("expected Remove to return true")
	}
	if w.Contains("billing") {
		t.Error("entry should be gone after Remove")
	}
}

func TestRemove_MissingEntry_ReturnsFalse(t *testing.T) {
	w := watchlist.New()
	if w.Remove("ghost") {
		t.Error("expected Remove to return false for unknown service")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(makeEntry("svc-a"))
	_ = w.Add(makeEntry("svc-b"))
	all := w.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestFilterByTag_MatchingEntries(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(makeEntry("api", "prod", "critical"))
	_ = w.Add(makeEntry("worker", "staging"))
	_ = w.Add(makeEntry("cron", "prod"))

	results := w.FilterByTag("prod")
	if len(results) != 2 {
		t.Fatalf("expected 2 prod entries, got %d", len(results))
	}
}

func TestFilterByTag_NoTags_ReturnsAll(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(makeEntry("svc-a", "prod"))
	_ = w.Add(makeEntry("svc-b", "staging"))

	results := w.FilterByTag()
	if len(results) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(results))
	}
}

func TestFilterByTag_CaseInsensitive(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(makeEntry("svc", "PROD"))
	results := w.FilterByTag("prod")
	if len(results) != 1 {
		t.Fatalf("expected 1 result for case-insensitive tag match, got %d", len(results))
	}
}

func TestFilterByTag_NoMatch_ReturnsEmpty(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(makeEntry("svc", "staging"))
	results := w.FilterByTag("prod")
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
