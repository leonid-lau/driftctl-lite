package tfstate

import "strings"

// DefaultCategoryFilterOptions returns options with case-insensitive matching enabled.
func DefaultCategoryFilterOptions() FilterOptions {
	return FilterOptions{CaseInsensitive: true}
}

// FilterByCategory returns resources whose "category" attribute matches the given value.
// If category is empty, all resources are returned.
func FilterByCategory(resources []Resource, category string, opts FilterOptions) []Resource {
	if category == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchCategory(r, category, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByCategories returns resources matching ANY of the provided categories (OR semantics).
// If the slice is empty, all resources are returned.
func FilterByCategories(resources []Resource, categories []string, opts FilterOptions) []Resource {
	if len(categories) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, c := range categories {
			if matchCategory(r, c, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchCategory(r Resource, category string, opts FilterOptions) bool {
	v, ok := r.Attributes["category"]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	if opts.CaseInsensitive {
		return strings.EqualFold(s, category)
	}
	return s == category
}
