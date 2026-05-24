// Package threshold classifies drift severity based on the number of
// fields that have drifted from their declared state.
//
// # Levels
//
// Four levels are defined:
//
//	none   – no drift detected
//	low    – a small number of fields differ
//	medium – a moderate number of fields differ
//	high   – many fields differ; immediate attention recommended
//
// # Usage
//
//	cfg := threshold.DefaultConfig()
//	eval, err := threshold.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	level := eval.Classify(driftedFieldCount)
//
// ClassifyReport is a convenience wrapper that operates on a full
// drift.Report and returns a Result slice.
package threshold
