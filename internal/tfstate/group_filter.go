package tfstate

import "strings"

// GroupFilterOptions controls matching behaviour for group filters.
type GroupFilterOptions struct {
	CaseSensitive bool
}

// DefaultGroupFilterOptions returns the default options (case-insensitive).
func DefaultGroupFilterOptions() GroupFilterOptions {
	return GroupFilterOptions{CaseSensitive: false}
}

// FilterByGroup returns resources whose "group" attribute matches the given
// value. An empty group value returns all resources unchanged.
func FilterByGroup(resources []Resource, group string, opts GroupFilterOptions) []Resource {
	if group == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchGroup(r, group, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByGroups returns resources that match ANY of the provided group values
// (OR semantics). An empty slice returns all resources unchanged.
func FilterByGroups(resources []Resource, groups []string, opts GroupFilterOptions) []Resource {
	if len(groups) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, g := range groups {
			if matchGroup(r, g, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchGroup(r Resource, group string, opts GroupFilterOptions) bool {
	v, ok := r.Attributes["group"]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	if opts.CaseSensitive {
		return s == group
	}
	return strings.EqualFold(s, group)
}
