package sampler

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Rate != 1.0 {
		t.Fatalf("expected default rate 1.0, got %v", cfg.Rate)
	}
}

func TestNew_InvalidRate_Zero(t *testing.T) {
	_, err := New(Config{Rate: 0.0})
	if err == nil {
		t.Fatal("expected error for rate=0.0, got nil")
	}
}

func TestNew_InvalidRate_Negative(t *testing.T) {
	_, err := New(Config{Rate: -0.5})
	if err == nil {
		t.Fatal("expected error for negative rate, got nil")
	}
}

func TestNew_InvalidRate_AboveOne(t *testing.T) {
	_, err := New(Config{Rate: 1.1})
	if err == nil {
		t.Fatal("expected error for rate > 1.0, got nil")
	}
}

func TestNew_ValidRate(t *testing.T) {
	s, err := New(Config{Rate: 0.5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
}

func TestSample_RateOne_AlwaysTrue(t *testing.T) {
	s, _ := New(DefaultConfig())
	for i := 0; i < 20; i++ {
		if !s.Sample("svc") {
			t.Fatal("rate=1.0 should always sample")
		}
	}
}

func TestSample_RateVeryLow_AlwaysFalse(t *testing.T) {
	// inject a rand that always returns 0.9 — above any low rate
	s, _ := newWithRand(Config{Rate: 0.1}, func() float64 { return 0.95 })
	if s.Sample("svc") {
		t.Fatal("expected Sample to return false when rand > rate")
	}
}

func TestSample_RateHigh_AlwaysTrue(t *testing.T) {
	// inject a rand that always returns 0.1 — below rate of 0.9
	s, _ := newWithRand(Config{Rate: 0.9}, func() float64 { return 0.05 })
	if !s.Sample("svc") {
		t.Fatal("expected Sample to return true when rand < rate")
	}
}

func TestSampleAll_FiltersCorrectly(t *testing.T) {
	calls := 0
	sequence := []float64{0.1, 0.9, 0.2, 0.8} // alternating pass/fail
	s, _ := newWithRand(Config{Rate: 0.5}, func() float64 {
		v := sequence[calls%len(sequence)]
		calls++
		return v
	})

	services := []string{"alpha", "beta", "gamma", "delta"}
	result := s.SampleAll(services)

	// 0.1 < 0.5 → alpha in; 0.9 >= 0.5 → beta out; 0.2 < 0.5 → gamma in; 0.8 >= 0.5 → delta out
	if len(result) != 2 {
		t.Fatalf("expected 2 sampled services, got %d: %v", len(result), result)
	}
	if result[0] != "alpha" || result[1] != "gamma" {
		t.Fatalf("unexpected sampled services: %v", result)
	}
}

func TestSampleAll_EmptyInput(t *testing.T) {
	s, _ := New(DefaultConfig())
	result := s.SampleAll([]string{})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}
