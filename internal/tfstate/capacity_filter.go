package tfstate

import "strings"

// CapacityFilterOptions controls matching behaviour for capacity filtering.
type CapacityFilterOptions struct {
	CaseSensitive bool
}

// DefaultCapacityFilterOptions returns the default options (case-insensitive).
func DefaultCapacityFilterOptions() CapacityFilterOptions {
	return CapacityFilterOptions{CaseSensitive: false}
}

// FilterByCapacity returns resources whose "capacity" attribute matches the
// given value. An empty capacity string returns all resources unchanged.
func FilterByCapacity(resources []Resource, capacity string, opts CapacityFilterOptions) []Resource {
	if capacity == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchCapacity(r, capacity, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByCapacities returns resources that match ANY of the given capacity
// values (OR semantics). An empty slice returns all resources unchanged.
func FilterByCapacities(resources []Resource, capacities []string, opts CapacityFilterOptions) []Resource {
	if len(capacities) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, c := range capacities {
			if matchCapacity(r, c, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchCapacity(r Resource, capacity string, opts CapacityFilterOptions) bool {
	v, ok := r.Attributes["capacity"]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	if opts.CaseSensitive {
		return s == capacity
	}
	return strings.EqualFold(s, capacity)
}
