package tfstate

import "strings"

// DefaultRuntimeFilterOptions returns case-insensitive options.
func DefaultRuntimeFilterOptions() FilterOptions {
	return FilterOptions{CaseSensitive: false}
}

// FilterByRuntime returns resources whose "runtime" attribute matches the given value.
// If runtime is empty, all resources are returned.
func FilterByRuntime(resources []Resource, runtime string, opts FilterOptions) []Resource {
	if runtime == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		val, ok := r.Attributes["runtime"]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		if matchRuntime(s, runtime, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByRuntimes returns resources matching any of the given runtime values (OR semantics).
func FilterByRuntimes(resources []Resource, runtimes []string, opts FilterOptions) []Resource {
	if len(runtimes) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		val, ok := r.Attributes["runtime"]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		for _, rt := range runtimes {
			if matchRuntime(s, rt, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchRuntime(val, target string, opts FilterOptions) bool {
	if opts.CaseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
