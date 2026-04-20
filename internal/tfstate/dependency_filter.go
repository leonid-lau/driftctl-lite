package tfstate

import "strings"

// DefaultDependencyFilterOptions returns options with case-insensitive matching.
func DefaultDependencyFilterOptions() FilterOptions {
	return FilterOptions{CaseSensitive: false}
}

// FilterOptions controls matching behaviour for dependency filters.
type FilterOptions struct {
	CaseSensitive bool
}

// FilterByDependency returns resources whose "dependency" attribute matches dep.
// If dep is empty, all resources are returned unchanged.
func FilterByDependency(resources []Resource, dep string, opts FilterOptions) []Resource {
	if dep == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["dependency"]
		if !ok {
			continue
		}
		s, _ := v.(string)
		if matchDep(s, dep, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByDependencies returns resources that match ANY of the supplied deps
// (OR semantics).
func FilterByDependencies(resources []Resource, deps []string, opts FilterOptions) []Resource {
	if len(deps) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["dependency"]
		if !ok {
			continue
		}
		s, _ := v.(string)
		for _, d := range deps {
			if matchDep(s, d, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchDep(val, target string, opts FilterOptions) bool {
	if opts.CaseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
