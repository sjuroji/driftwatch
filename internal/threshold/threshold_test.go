package threshold_test

import (
	"testing"

	"github.com/driftwatch/internal/threshold"
)

func TestDefaultConfig(t *testing.T) {
	cfg := threshold.DefaultConfig()
	if cfg.LowAt <= 0 {
		t.Fatalf("LowAt must be positive, got %d", cfg.LowAt)
	}
	if cfg.MediumAt <= cfg.LowAt {
		t.Fatalf("MediumAt(%d) must be > LowAt(%d)", cfg.MediumAt, cfg.LowAt)
	}
	if cfg.HighAt <= cfg.MediumAt {
		t.Fatalf("HighAt(%d) must be > MediumAt(%d)", cfg.HighAt, cfg.MediumAt)
	}
}

func TestNew_InvalidBoundaries(t *testing.T) {
	cases := []struct {
		name string
		cfg  threshold.Config
	}{
		{"zero low", threshold.Config{LowAt: 0, MediumAt: 3, HighAt: 6}},
		{"medium <= low", threshold.Config{LowAt: 3, MediumAt: 3, HighAt: 6}},
		{"high <= medium", threshold.Config{LowAt: 1, MediumAt: 3, HighAt: 3}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := threshold.New(tc.cfg)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestClassify_None(t *testing.T) {
	e, _ := threshold.New(threshold.DefaultConfig())
	if got := e.Classify(0); got != threshold.LevelNone {
		t.Fatalf("expected none, got %s", got)
	}
}

func TestClassify_Low(t *testing.T) {
	e, _ := threshold.New(threshold.DefaultConfig()) // low=1, medium=3, high=6
	for _, n := range []int{1, 2} {
		if got := e.Classify(n); got != threshold.LevelLow {
			t.Fatalf("Classify(%d): expected low, got %s", n, got)
		}
	}
}

func TestClassify_Medium(t *testing.T) {
	e, _ := threshold.New(threshold.DefaultConfig())
	for _, n := range []int{3, 4, 5} {
		if got := e.Classify(n); got != threshold.LevelMedium {
			t.Fatalf("Classify(%d): expected medium, got %s", n, got)
		}
	}
}

func TestClassify_High(t *testing.T) {
	e, _ := threshold.New(threshold.DefaultConfig())
	for _, n := range []int{6, 10, 100} {
		if got := e.Classify(n); got != threshold.LevelHigh {
			t.Fatalf("Classify(%d): expected high, got %s", n, got)
		}
	}
}

func TestClassify_NegativeIsNone(t *testing.T) {
	e, _ := threshold.New(threshold.DefaultConfig())
	if got := e.Classify(-5); got != threshold.LevelNone {
		t.Fatalf("expected none for negative input, got %s", got)
	}
}
