// Package sampler provides probabilistic sampling for drift detection runs.
//
// In environments with a large number of services, running a full drift check
// on every service at every tick may be expensive. The Sampler allows a
// configurable fraction of services to be selected for checking in each cycle,
// reducing load while maintaining broad coverage over time.
//
// Usage:
//
//	s, err := sampler.New(sampler.Config{Rate: 0.25})
//	if err != nil {
//		log.Fatal(err)
//	}
//	selected := s.SampleAll(allServiceNames)
//
// A Rate of 1.0 (the default) disables sampling — every service is included.
package sampler
