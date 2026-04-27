package tfstate

import "strings"

// PoolFilterOptions configures pool filtering behaviour.
type PoolFilterOptions struct {
	CaseSensitive bool
}

// DefaultPoolFilterOptions returns the default pool filter options.
func DefaultPoolFilterOptions() PoolFilterOptions {
	return PoolFilterOptions{CaseSensitive: false}
}

// FilterByPool returns resources whose "pool" attribute matches the given value.
// If pool is empty, all resources are returned.
func FilterByPool(resources []Resource, pool string, opts PoolFilterOptions) []Resource {
	if pool == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchPool(r, pool, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByPools returns resources that match ANY of the given pool values.
// If the pools slice is empty, all resources are returned.
func FilterByPools(resources []Resource, pools []string, opts PoolFilterOptions) []Resource {
	if len(pools) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, p := range pools {
			if matchPool(r, p, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchPool(r Resource, pool string, opts PoolFilterOptions) bool {
	v, ok := r.Attributes["pool"]
	if !ok {
		return false
	}
	if opts.CaseSensitive {
		return v == pool
	}
	return strings.EqualFold(v, pool)
}
