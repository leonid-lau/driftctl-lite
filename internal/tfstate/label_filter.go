package tfstate

// LabelFilterOptions controls how label-based filtering is applied.
type LabelFilterOptions struct {
	// RequireAll requires all labels to match (AND semantics).
	// When false, any match is sufficient (OR semantics).
	RequireAll bool
}

// DefaultLabelFilterOptions returns options with AND semantics.
func DefaultLabelFilterOptions() LabelFilterOptions {
	return LabelFilterOptions{RequireAll: true}
}

// FilterByLabel returns resources that have the given label key set to the given value.
// If value is empty, any resource with the key present is returned.
func FilterByLabel(resources []Resource, key, value string) []Resource {
	if key == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Labels[key]
		if !ok {
			continue
		}
		if value == "" || v == value {
			out = append(out, r)
		}
	}
	return out
}

// FilterByLabels returns resources matching all or any of the given labels
// depending on opts.RequireAll.
func FilterByLabels(resources []Resource, labels map[string]string, opts LabelFilterOptions) []Resource {
	if len(labels) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		matched := 0
		for k, v := range labels {
			if rv, ok := r.Labels[k]; ok && (v == "" || rv == v) {
				matched++
			}
		}
		if opts.RequireAll && matched == len(labels) {
			out = append(out, r)
		} else if !opts.RequireAll && matched > 0 {
			out = append(out, r)
		}
	}
	return out
}
