package tag_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/tag"
)

// taggedItem is a test helper that implements tag.Tagged.
type taggedItem struct {
	name string
	tags []string
}

func (ti taggedItem) Tags() []string { return ti.tags }

func TestNewSet_NormalisesCase(t *testing.T) {
	s := tag.NewSet([]string{"Prod", "EU-WEST", "  auth  "})
	for _, want := range []string{"prod", "eu-west", "auth"} {
		if !s.Contains(want) {
			t.Errorf("expected set to contain %q", want)
		}
	}
}

func TestSet_Contains_CaseInsensitive(t *testing.T) {
	s := tag.NewSet([]string{"critical"})
	if !s.Contains("CRITICAL") {
		t.Error("Contains should be case-insensitive")
	}
	if s.Contains("warning") {
		t.Error("Contains should return false for absent tag")
	}
}

func TestSet_Intersects_True(t *testing.T) {
	a := tag.NewSet([]string{"prod", "eu"})
	b := tag.NewSet([]string{"eu", "us"})
	if !a.Intersects(b) {
		t.Error("expected sets to intersect on 'eu'")
	}
}

func TestSet_Intersects_False(t *testing.T) {
	a := tag.NewSet([]string{"prod"})
	b := tag.NewSet([]string{"staging"})
	if a.Intersects(b) {
		t.Error("expected sets not to intersect")
	}
}

func TestSet_Slice_Sorted(t *testing.T) {
	s := tag.NewSet([]string{"zebra", "apple", "mango"})
	got := s.Slice()
	want := []string{"apple", "mango", "zebra"}
	if len(got) != len(want) {
		t.Fatalf("expected %d tags, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: want %q, got %q", i, want[i], got[i])
		}
	}
}

func TestMatch_EmptyRequired_ReturnsAll(t *testing.T) {
	items := []tag.Tagged{
		taggedItem{name: "a", tags: []string{"prod"}},
		taggedItem{name: "b", tags: []string{"staging"}},
	}
	got := tag.Match(tag.NewSet(nil), items)
	if len(got) != 2 {
		t.Errorf("expected 2 items, got %d", len(got))
	}
}

func TestMatch_FiltersByTag(t *testing.T) {
	items := []tag.Tagged{
		taggedItem{name: "svc-a", tags: []string{"prod", "eu"}},
		taggedItem{name: "svc-b", tags: []string{"staging"}},
		taggedItem{name: "svc-c", tags: []string{"prod", "us"}},
	}
	required := tag.NewSet([]string{"prod"})
	got := tag.Match(required, items)
	if len(got) != 2 {
		t.Errorf("expected 2 matches, got %d", len(got))
	}
}

func TestMatch_NoMatches_ReturnsNil(t *testing.T) {
	items := []tag.Tagged{
		taggedItem{name: "svc-a", tags: []string{"staging"}},
	}
	got := tag.Match(tag.NewSet([]string{"prod"}), items)
	if len(got) != 0 {
		t.Errorf("expected 0 matches, got %d", len(got))
	}
}
