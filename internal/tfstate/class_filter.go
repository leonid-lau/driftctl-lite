package tfstate

import "strings"

// ClassFilterOptions configures FilterByClass behaviour.
type ClassFilterOptions struct {
	CaseSensitive bool
}

// DefaultClassFilterOptions returns the default options (case-insensitive).
func DefaultClassFilterOptions() ClassFilterOptions {
	return ClassFilterOptions{CaseSensitive: false}
}

// FilterByClass returns resources whose "class" attribute matches the given
// value. An empty class string returns all resources unchanged.
func FilterByClass(resources []Resource, class string, opts ClassFilterOptions) []Resource {
	if class == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchClass(r, class, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByClasses returns resources that match ANY of the provided class
// values (OR semantics). An empty slice returns all resources unchanged.
func FilterByClasses(resources []Resource, classes []string, opts ClassFilterOptions) []Resource {
	if len(classes) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, c := range classes {
			if matchClass(r, c, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchClass(r Resource, class string, opts ClassFilterOptions) bool {
	v, ok := r.Attributes["class"]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	if opts.CaseSensitive {
		return s == class
	}
	return strings.EqualFold(s, class)
}
