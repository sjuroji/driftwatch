// Package audit provides structured, append-only audit logging for
// driftwatch operations.
//
// Each audit event is written as a single JSON line containing a timestamp,
// severity level, human-readable message, optional service name, and an
// arbitrary key/value fields map.
//
// Usage:
//
//	logger := audit.New(os.Stderr)
//	logger.Info("drift scan started", "", nil)
//	logger.Warn("replica drift detected", "auth-service", map[string]string{
//		"expected": "3",
//		"actual":   "1",
//	})
//
// The output is suitable for ingestion by log aggregation systems such as
// Loki, Elasticsearch, or CloudWatch Logs.
package audit
