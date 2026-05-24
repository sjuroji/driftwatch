// Package trend provides drift-trend analysis over a series of historical
// drift reports.
//
// Given a slice of drift.Report values ordered oldest-first, Analyse
// computes a per-service drift rate and classifies each service as
// stable, worsening, or improving by comparing the drift rate of the
// earlier half of the window to the later half.
//
// Typical usage:
//
//	reports, err := historyStore.List()
//	if err != nil { ... }
//
//	tr, err := trend.Analyse(reports, 24*time.Hour)
//	if err != nil { ... }
//
//	for _, svc := range tr.Services {
//		fmt.Printf("%s drift=%.0f%% direction=%s\n",
//			svc.Service, svc.DriftRate*100, svc.Direction)
//	}
package trend
