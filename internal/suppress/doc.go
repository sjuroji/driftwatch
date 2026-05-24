// Package suppress manages drift suppression rules for driftwatch.
//
// A suppression rule tells driftwatch to ignore known, accepted deviations
// between a deployed service and its declared manifest. Rules can target an
// entire service or a specific field (e.g. "replicas" managed by an
// autoscaler), and may carry an optional expiry time after which the
// suppression is automatically lifted.
//
// Example usage:
//
//	rules := []suppress.Rule{
//		{
//			Service:   "auth-service",
//			Field:     "replicas",
//			Reason:    "managed by HPA",
//			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
//		},
//	}
//	s := suppress.New(rules, nil)
//	if s.IsSuppressed("auth-service", "replicas") {
//		// skip reporting this drift
//	}
package suppress
