package tfstate

import "strings"

// DefaultProviderFilterOptions returns case-insensitive matching options.
type ProviderFilterOptions struct {
	CaseSensitive bool
}

func DefaultProviderFilterOptions() ProviderFilterOptions {
	return ProviderFilterOptions{CaseSensitive: false}
}

// FilterByProvider returns resources whose Provider field matches the given value.
// An empty provider string returns all resources unchanged.
func FilterByProvider(resources []Resource, provider string, opts ProviderFilterOptions) []Resource {
	if provider == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		rp := r.Provider
		p := provider
		if !opts.CaseSensitive {
			rp = strings.ToLower(rp)
			p = strings.ToLower(p)
		}
		if rp == p {
			out = append(out, r)
		}
	}
	return out
}

// FilterByProviders returns resources matching any of the given provider values (OR semantics).
func FilterByProviders(resources []Resource, providers []string, opts ProviderFilterOptions) []Resource {
	if len(providers) == 0 {
		return resources
	}
	set := make(map[string]struct{}, len(providers))
	for _, p := range providers {
		key := p
		if !opts.CaseSensitive {
			key = strings.ToLower(p)
		}
		set[key] = struct{}{}
	}
	var out []Resource
	for _, r := range resources {
		rp := r.Provider
		if !opts.CaseSensitive {
			rp = strings.ToLower(rp)
		}
		if _, ok := set[rp]; ok {
			out = append(out, r)
		}
	}
	return out
}
