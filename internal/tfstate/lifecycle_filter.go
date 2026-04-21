package tfstate

import "strings"

// DefaultLifecycleFilterOptions returns options with case-insensitive matching.
type LifecycleFilterOptions struct {
	CaseSensitive bool
}

func DefaultLifecycleFilterOptions() LifecycleFilterOptions {
	return LifecycleFilterOptions{CaseSensitive: false}
}

// FilterByLifecycle returns resources whose "lifecycle" attribute matches the given value.
// An empty lifecycle value returns all resources unchanged.
func FilterByLifecycle(resources []Resource, lifecycle string, opts LifecycleFilterOptions) []Resource {
	if lifecycle == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["lifecycle"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		if matchLifecycle(s, lifecycle, opts.CaseSensitive) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByLifecycles returns resources matching any of the provided lifecycle values (OR semantics).
func FilterByLifecycles(resources []Resource, lifecycles []string, opts LifecycleFilterOptions) []Resource {
	if len(lifecycles) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["lifecycle"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		for _, lc := range lifecycles {
			if matchLifecycle(s, lc, opts.CaseSensitive) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchLifecycle(val, target string, caseSensitive bool) bool {
	if caseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
