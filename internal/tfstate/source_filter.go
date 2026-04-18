package tfstate

import "strings"

// DefaultSourceFilterOptions returns case-insensitive matching by default.
type SourceFilterOptions struct {
	CaseSensitive bool
}

func DefaultSourceFilterOptions() SourceFilterOptions {
	return SourceFilterOptions{CaseSensitive: false}
}

// FilterBySource returns resources whose Source field matches the given value.
// An empty source returns all resources.
func FilterBySource(resources []Resource, source string, opts SourceFilterOptions) []Resource {
	if source == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		rv := r.Source
		sv := source
		if !opts.CaseSensitive {
			rv = strings.ToLower(rv)
			sv = strings.ToLower(sv)
		}
		if rv == sv {
			out = append(out, r)
		}
	}
	return out
}

// FilterBySources returns resources matching any of the given sources (OR semantics).
func FilterBySources(resources []Resource, sources []string, opts SourceFilterOptions) []Resource {
	if len(sources) == 0 {
		return resources
	}
	set := make(map[string]struct{}, len(sources))
	for _, s := range sources {
		key := s
		if !opts.CaseSensitive {
			key = strings.ToLower(key)
		}
		set[key] = struct{}{}
	}
	var out []Resource
	for _, r := range resources {
		rv := r.Source
		if !opts.CaseSensitive {
			rv = strings.ToLower(rv)
		}
		if _, ok := set[rv]; ok {
			out = append(out, r)
		}
	}
	return out
}
