package tfstate

import "strings"

// DefaultCreatedFilterOptions returns a CreatedFilterOptions with case-insensitive matching.
type CreatedFilterOptions struct {
	CaseSensitive bool
}

// DefaultCreatedFilterOptions returns the default options for filtering by created-by.
func DefaultCreatedFilterOptions() CreatedFilterOptions {
	return CreatedFilterOptions{CaseSensitive: false}
}

// FilterByCreated returns resources whose "created_by" attribute matches the given value.
// If createdBy is empty, all resources are returned.
func FilterByCreated(resources []Resource, createdBy string, opts CreatedFilterOptions) []Resource {
	if createdBy == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		val, ok := r.Attributes["created_by"]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		if matchCreated(s, createdBy, opts.CaseSensitive) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByCreatedValues returns resources matching any of the provided created-by values.
func FilterByCreatedValues(resources []Resource, values []string, opts CreatedFilterOptions) []Resource {
	if len(values) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		val, ok := r.Attributes["created_by"]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		for _, v := range values {
			if matchCreated(s, v, opts.CaseSensitive) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchCreated(attr, target string, caseSensitive bool) bool {
	if caseSensitive {
		return attr == target
	}
	return strings.EqualFold(attr, target)
}
