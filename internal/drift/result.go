package drift

// Status constants used by Result and consumers such as groupby.
const (
	StatusInSync  = "in-sync"
	StatusDrifted = "drifted"
	StatusUnknown = "unknown"
)

// Result holds the outcome of comparing a single service's declared manifest
// against its live state.
type Result struct {
	// Service is the canonical name of the service.
	Service string

	// Status is one of StatusInSync, StatusDrifted, or StatusUnknown.
	Status string

	// Fields lists the individual field names that differ between the
	// manifest and the live state. Empty when Status is StatusInSync.
	Fields []string

	// Labels carries arbitrary key/value metadata from the manifest,
	// enabling downstream grouping and filtering.
	Labels map[string]string
}
