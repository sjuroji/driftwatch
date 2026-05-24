package label_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/label"
)

func TestParseSelector_Empty(t *testing.T) {
	sel, err := label.ParseSelector("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sel.Empty() {
		t.Error("expected empty selector")
	}
}

func TestParseSelector_Valid(t *testing.T) {
	sel, err := label.ParseSelector("env=prod,team=platform")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sel.Empty() {
		t.Error("expected non-empty selector")
	}
}

func TestParseSelector_InvalidSegment(t *testing.T) {
	_, err := label.ParseSelector("env")
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestParseSelector_EmptyKey(t *testing.T) {
	_, err := label.ParseSelector("=value")
	if err == nil {
		t.Error("expected error for empty key")
	}
}

func TestMatches_AllPresent(t *testing.T) {
	sel, _ := label.ParseSelector("env=prod,team=platform")
	labels := map[string]string{"env": "prod", "team": "platform", "region": "us-east"}
	if !sel.Matches(labels) {
		t.Error("expected selector to match")
	}
}

func TestMatches_MissingKey(t *testing.T) {
	sel, _ := label.ParseSelector("env=prod")
	if sel.Matches(map[string]string{"team": "platform"}) {
		t.Error("expected selector not to match")
	}
}

func TestMatches_WrongValue(t *testing.T) {
	sel, _ := label.ParseSelector("env=prod")
	if sel.Matches(map[string]string{"env": "staging"}) {
		t.Error("expected selector not to match on wrong value")
	}
}

func TestMatches_CaseInsensitiveKey(t *testing.T) {
	sel, _ := label.ParseSelector("ENV=prod")
	if !sel.Matches(map[string]string{"env": "prod"}) {
		t.Error("expected case-insensitive key match")
	}
}

func TestMatches_EmptySelector_AlwaysTrue(t *testing.T) {
	sel, _ := label.ParseSelector("")
	if !sel.Matches(map[string]string{"env": "prod"}) {
		t.Error("empty selector should always match")
	}
	if !sel.Matches(nil) {
		t.Error("empty selector should match nil labels")
	}
}

func TestString_RoundTrip(t *testing.T) {
	sel, _ := label.ParseSelector("env=prod")
	if sel.String() == "" {
		t.Error("expected non-empty string representation")
	}
}
