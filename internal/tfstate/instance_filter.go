package tfstate

import "strings"

// InstanceFilterOptions controls matching behaviour for FilterByInstance.
type InstanceFilterOptions struct {
	CaseSensitive bool
}

// DefaultInstanceFilterOptions returns the default options (case-insensitive).
func DefaultInstanceFilterOptions() InstanceFilterOptions {
	return InstanceFilterOptions{CaseSensitive: false}
}

// FilterByInstance returns resources whose "instance" attribute matches the
// given value. An empty value returns all resources unchanged.
func FilterByInstance(resources []Resource, instance string, opts InstanceFilterOptions) []Resource {
	if instance == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchInstance(r, instance, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByInstances returns resources that match ANY of the provided instance
// values (OR semantics). An empty slice returns all resources unchanged.
func FilterByInstances(resources []Resource, instances []string, opts InstanceFilterOptions) []Resource {
	if len(instances) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, inst := range instances {
			if matchInstance(r, inst, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchInstance(r Resource, instance string, opts InstanceFilterOptions) bool {
	val, ok := r.Attributes["instance"]
	if !ok {
		return false
	}
	s, ok := val.(string)
	if !ok {
		return false
	}
	if opts.CaseSensitive {
		return s == instance
	}
	return strings.EqualFold(s, instance)
}
