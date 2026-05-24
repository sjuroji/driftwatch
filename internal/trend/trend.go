// Package trend analyses drift report history to identify services
// that are repeatedly or increasingly drifting over time.
package trend

import (
	"errors"
	"sort"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Direction describes the direction of a drift trend.
type Direction string

const (
	DirectionStable    Direction = "stable"
	DirectionWorsening Direction = "worsening"
	DirectionImproving Direction = "improving"
)

// ServiceTrend holds trend information for a single service.
type ServiceTrend struct {
	Service    string
	DriftCount int
	Total      int
	DriftRate  float64
	Direction  Direction
}

// Report is the output of Analyse.
type Report struct {
	Window   time.Duration
	Services []ServiceTrend
	GeneratedAt time.Time
}

// Analyse computes per-service drift trends from an ordered slice of
// drift reports. Reports must be sorted oldest-first. window is used
// only for documentation in the returned Report.
func Analyse(reports []drift.Report, window time.Duration) (Report, error) {
	if len(reports) == 0 {
		return Report{}, errors.New("trend: no reports provided")
	}

	type counts struct {
		drifted int
		total   int
	}

	serviceMap := make(map[string]*counts)

	for _, r := range reports {
		for _, entry := range r.Entries {
			c, ok := serviceMap[entry.Service]
			if !ok {
				c = &counts{}
				serviceMap[entry.Service] = c
			}
			c.total++
			if entry.Status == drift.StatusDrifted {
				c.drifted++
			}
		}
	}

	services := make([]ServiceTrend, 0, len(serviceMap))
	for svc, c := range serviceMap {
		rate := float64(c.drifted) / float64(c.total)
		services = append(services, ServiceTrend{
			Service:   svc,
			DriftCount: c.drifted,
			Total:     c.total,
			DriftRate: rate,
			Direction: direction(reports, svc),
		})
	}

	sort.Slice(services, func(i, j int) bool {
		return services[i].DriftRate > services[j].DriftRate
	})

	return Report{
		Window:      window,
		Services:    services,
		GeneratedAt: time.Now(),
	}, nil
}

// direction determines whether a service is getting better or worse by
// comparing the drift rate of the first half of reports to the second half.
func direction(reports []drift.Report, service string) Direction {
	mid := len(reports) / 2
	if mid == 0 {
		return DirectionStable
	}

	early := driftRate(reports[:mid], service)
	late := driftRate(reports[mid:], service)

	switch {
	case late > early:
		return DirectionWorsening
	case late < early:
		return DirectionImproving
	default:
		return DirectionStable
	}
}

func driftRate(reports []drift.Report, service string) float64 {
	var drifted, total int
	for _, r := range reports {
		for _, e := range r.Entries {
			if e.Service == service {
				total++
				if e.Status == drift.StatusDrifted {
					drifted++
				}
			}
		}
	}
	if total == 0 {
		return 0
	}
	return float64(drifted) / float64(total)
}
