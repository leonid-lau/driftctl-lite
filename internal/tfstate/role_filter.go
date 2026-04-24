package tfstate

import "strings"

// DefaultRoleFilterOptions returns case-insensitive matching options.
func DefaultRoleFilterOptions() FilterOptions {
	return FilterOptions{CaseSensitive: false}
}

// FilterOptions controls matching behaviour for role filters.
type FilterOptions struct {
	CaseSensitive bool
}

// FilterByRole returns resources whose "role" attribute matches the given value.
// An empty role returns all resources unchanged.
func FilterByRole(resources []Resource, role string, opts FilterOptions) []Resource {
	if role == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchRole(r, role, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByRoles returns resources that match ANY of the provided role values (OR semantics).
func FilterByRoles(resources []Resource, roles []string, opts FilterOptions) []Resource {
	if len(roles) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, role := range roles {
			if matchRole(r, role, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchRole(r Resource, role string, opts FilterOptions) bool {
	val, ok := r.Attributes["role"]
	if !ok {
		return false
	}
	s, ok := val.(string)
	if !ok {
		return false
	}
	if opts.CaseSensitive {
		return s == role
	}
	return strings.EqualFold(s, role)
}
