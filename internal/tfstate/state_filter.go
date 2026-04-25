package tfstate

import "strings"

// StateFilterOptions controls matching behaviour for FilterByState.
type StateFilterOptions struct {
	// CaseSensitive disables the default case-insensitive comparison.
	CaseSensitive bool
}

// DefaultStateFilterOptions returns the recommended defaults.
func DefaultStateFilterOptions() StateFilterOptions {
	return StateFilterOptions{CaseSensitive: false}
}

// FilterByState returns resources whose "state" attribute matches target.
// An empty target returns all resources unchanged.
func FilterByState(resources []Resource, target string, opts StateFilterOptions) []Resource {
	if target == "" {
		return resources
	}
	cmp := target
	if !opts.CaseSensitive {
		cmp = strings.ToLower(target)
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["state"]
		if !ok {
			continue
		}
		v2, _ := v.(string)
		candidate := v2
		if !opts.CaseSensitive {
			candidate = strings.ToLower(v2)
		}
		if candidate == cmp {
			out = append(out, r)
		}
	}
	return out
}

// FilterByStates returns resources whose "state" attribute matches any of
// the provided targets (OR semantics). An empty slice returns all resources.
func FilterByStates(resources []Resource, targets []string, opts StateFilterOptions) []Resource {
	if len(targets) == 0 {
		return resources
	}
	set := make(map[string]struct{}, len(targets))
	for _, t := range targets {
		key := t
		if !opts.CaseSensitive {
			key = strings.ToLower(t)
		}
		set[key] = struct{}{}
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["state"]
		if !ok {
			continue
		}
		v2, _ := v.(string)
		candidate := v2
		if !opts.CaseSensitive {
			candidate = strings.ToLower(v2)
		}
		if _, found := set[candidate]; found {
			out = append(out, r)
		}
	}
	return out
}
