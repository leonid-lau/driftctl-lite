package tfstate

import "strings"

// DefaultComponentFilterOptions returns options with case-insensitive matching.
type ComponentFilterOptions struct {
	CaseSensitive bool
}

func DefaultComponentFilterOptions() ComponentFilterOptions {
	return ComponentFilterOptions{CaseSensitive: false}
}

// FilterByComponent returns resources whose "component" attribute matches value.
func FilterByComponent(resources []Resource, component string, opts ComponentFilterOptions) []Resource {
	if component == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchComponent(r, component, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByComponents returns resources matching any of the given components (OR semantics).
func FilterByComponents(resources []Resource, components []string, opts ComponentFilterOptions) []Resource {
	if len(components) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, c := range components {
			if matchComponent(r, c, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchComponent(r Resource, component string, opts ComponentFilterOptions) bool {
	v, ok := r.Attributes["component"]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	if opts.CaseSensitive {
		return s == component
	}
	return strings.EqualFold(s, component)
}
