package label

// Labelled is any type that exposes a Labels map.
type Labelled interface {
	GetLabels() map[string]string
}

// Filter returns only those items from the slice whose labels satisfy the selector.
// If the selector is empty, all items are returned unchanged.
func Filter[T Labelled](items []T, sel Selector) []T {
	if sel.Empty() {
		return items
	}
	out := make([]T, 0, len(items))
	for _, item := range items {
		if sel.Matches(item.GetLabels()) {
			out = append(out, item)
		}
	}
	return out
}

// GroupBy partitions items into buckets keyed by the value of labelKey.
// Items that do not carry the label are placed under the empty-string key.
func GroupBy[T Labelled](items []T, labelKey string) map[string][]T {
	result := make(map[string][]T)
	for _, item := range items {
		v := item.GetLabels()[labelKey]
		result[v] = append(result[v], item)
	}
	return result
}
