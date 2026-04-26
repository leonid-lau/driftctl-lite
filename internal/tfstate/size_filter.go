package tfstate

import "strings"

// SizeFilterOptions controls matching behaviour for FilterBySize.
type SizeFilterOptions struct {
	// CaseSensitive disables case-folding when true.
	CaseSensitive bool
}

// DefaultSizeFilterOptions returns the default options (case-insensitive).
func DefaultSizeFilterOptions() SizeFilterOptions {
	return SizeFilterOptions{CaseSensitive: false}
}

// FilterBySize returns resources whose "size" attribute matches the given
// value. An empty size string returns all resources unchanged.
func FilterBySize(resources []Resource, size string, opts SizeFilterOptions) []Resource {
	if size == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["size"].(string)
		if matchSize(v, size, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterBySizes returns resources whose "size" attribute matches ANY of the
// provided values (OR semantics).
func FilterBySizes(resources []Resource, sizes []string, opts SizeFilterOptions) []Resource {
	if len(sizes) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["size"].(string)
		for _, s := range sizes {
			if matchSize(v, s, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchSize(val, want string, opts SizeFilterOptions) bool {
	if opts.CaseSensitive {
		return val == want
	}
	return strings.EqualFold(val, want)
}
