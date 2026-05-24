// Package sampler provides probabilistic sampling for drift reports,
// allowing high-volume environments to reduce noise by processing only
// a fraction of checks while preserving statistical significance.
package sampler

import (
	"errors"
	"math/rand"
)

// Config holds configuration for the Sampler.
type Config struct {
	// Rate is the fraction of items to sample, in the range (0.0, 1.0].
	// A rate of 1.0 means every item is sampled; 0.5 means ~50%.
	Rate float64
}

// DefaultConfig returns a Config that samples every item.
func DefaultConfig() Config {
	return Config{Rate: 1.0}
}

// Sampler decides whether a given service check should be processed.
type Sampler struct {
	cfg  Config
	rand func() float64
}

// New creates a new Sampler with the given Config.
// Returns an error if Rate is not in (0.0, 1.0].
func New(cfg Config) (*Sampler, error) {
	if cfg.Rate <= 0.0 || cfg.Rate > 1.0 {
		return nil, errors.New("sampler: rate must be in the range (0.0, 1.0]")
	}
	return &Sampler{
		cfg:  cfg,
		rand: rand.Float64,
	}, nil
}

// newWithRand creates a Sampler with an injected random source for testing.
func newWithRand(cfg Config, randFn func() float64) (*Sampler, error) {
	s, err := New(cfg)
	if err != nil {
		return nil, err
	}
	s.rand = randFn
	return s, nil
}

// Sample returns true if the named service should be processed in this
// check cycle. The decision is probabilistic based on the configured rate.
func (s *Sampler) Sample(service string) bool {
	_ = service // reserved for future deterministic/hash-based sampling
	return s.rand() < s.cfg.Rate
}

// SampleAll filters a slice of service names, returning only those
// that pass the sampling decision.
func (s *Sampler) SampleAll(services []string) []string {
	out := make([]string, 0, len(services))
	for _, svc := range services {
		if s.Sample(svc) {
			out = append(out, svc)
		}
	}
	return out
}
