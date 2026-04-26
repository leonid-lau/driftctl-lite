package tfstate

import "strings"

// DefaultFormatFilterOptions returns options with case-insensitive matching enabled.
func DefaultFormatFilterOptions() FilterOptions {
	return FilterOptions{CaseInsensitive: true}
}

// FilterOptions holds configuration for filter operations.
type FilterOptions struct {
	CaseInsensitive bool
}

// FilterByFormat returns resources whose "format" attribute matches the given value.
// If format is empty, all resources are returned.
func FilterByFormat(resources []Resource, format string, opts FilterOptions) []Resource {
	if format == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["format"]
		if !ok {
			continue
		}
		if matchFormat(v, format, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByFormats returns resources matching any of the given format values (OR semantics).
func FilterByFormats(resources []Resource, formats []string, opts FilterOptions) []Resource {
	if len(formats) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["format"]
		if !ok {
			continue
		}
		for _, f := range formats {
			if matchFormat(v, f, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchFormat(attr, target string, opts FilterOptions) bool {
	if opts.CaseInsensitive {
		return strings.EqualFold(attr, target)
	}
	return attr == target
}
