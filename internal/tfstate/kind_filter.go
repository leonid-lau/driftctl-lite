package tfstate

import "strings"

// DefaultKindFilterOptions returns case-insensitive matching options.
type KindFilterOptions struct {
	CaseSensitive bool
}

func DefaultKindFilterOptions() KindFilterOptions {
	return KindFilterOptions{CaseSensitive: false}
}

// FilterByKind returns resources whose Kind matches the given value.
// An empty kind returns all resources.
func FilterByKind(resources []Resource, kind string, opts KindFilterOptions) []Resource {
	if kind == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		rk := r.Kind
		k := kind
		if !opts.CaseSensitive {
			rk = strings.ToLower(rk)
			k = strings.ToLower(k)
		}
		if rk == k {
			out = append(out, r)
		}
	}
	return out
}

// FilterByKinds returns resources matching any of the provided kinds (OR semantics).
func FilterByKinds(resources []Resource, kinds []string, opts KindFilterOptions) []Resource {
	if len(kinds) == 0 {
		return resources
	}
	set := make(map[string]struct{}, len(kinds))
	for _, k := range kinds {
		if !opts.CaseSensitive {
			k = strings.ToLower(k)
		}
		set[k] = struct{}{}
	}
	var out []Resource
	for _, r := range resources {
		rk := r.Kind
		if !opts.CaseSensitive {
			rk = strings.ToLower(rk)
		}
		if _, ok := set[rk]; ok {
			out = append(out, r)
		}
	}
	return out
}
