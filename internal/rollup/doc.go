// Package rollup provides aggregation of multiple drift.Report values into a
// single summarised view.
//
// Given a slice of reports produced over time (e.g. loaded from the history
// store), Aggregate computes per-service drift frequency metrics:
//
//	res, err := rollup.Aggregate(reports)
//	for _, s := range res.Services {
//		fmt.Printf("%s drifted in %.0f%% of runs\n", s.Service, s.DriftRate*100)
//	}
//
// Services are returned sorted by DriftRate descending so that the most
// problematic services appear first.
package rollup
