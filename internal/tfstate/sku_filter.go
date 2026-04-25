package tfstate

import "strings"

// DefaultSKUFilterOptions returns case-insensitive matching by default.
type SKUFilterOptions struct {
	CaseSensitive bool
}

func DefaultSKUFilterOptions() SKUFilterOptions {
	return SKUFilterOptions{CaseSensitive: false}
}

// FilterBySKU returns resources whose "sku" attribute matches the given value.
// If sku is empty, all resources are returned unchanged.
func FilterBySKU(resources []Resource, sku string, opts SKUFilterOptions) []Resource {
	if sku == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchSKU(r, sku, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterBySKUs returns resources matching ANY of the provided SKU values (OR semantics).
func FilterBySKUs(resources []Resource, skus []string, opts SKUFilterOptions) []Resource {
	if len(skus) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, s := range skus {
			if matchSKU(r, s, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchSKU(r Resource, sku string, opts SKUFilterOptions) bool {
	val, ok := r.Attributes["sku"]
	if !ok {
		return false
	}
	s, ok := val.(string)
	if !ok {
		return false
	}
	if opts.CaseSensitive {
		return s == sku
	}
	return strings.EqualFold(s, sku)
}
