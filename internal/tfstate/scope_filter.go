package tfstate

import "strings"

// ScopeFilterOptions controls matching behaviour for scope filters.
type ScopeFilterOptions struct {
	CaseSensitive bool
}

// DefaultScopeFilterOptions returns the default options (case-insensitive).
func DefaultScopeFilterOptions() ScopeFilterOptions {
	return ScopeFilterOptions{CaseSensitive: false}
}

// FilterByScope returns resources whose Attributes["scope"] matches the given value.
// An empty scope returns all resources unchanged.
func FilterByScope(resources []Resource, scope string, opts ScopeFilterOptions) []Resource {
	if scope == "" {
		return resources
	}
	want := scope
	if !opts.CaseSensitive {
		want = strings.ToLower(scope)
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["scope"]
		if !ok {
			continue
		}
		v2 := v
		if !opts.CaseSensitive {
			v2 = strings.ToLower(v)
		}
		if v2 == want {
			out = append(out, r)
		}
	}
	return out
}

// FilterByScopes returns resources matching ANY of the provided scopes (OR semantics).
func FilterByScopes(resources []Resource, scopes []string, opts ScopeFilterOptions) []Resource {
	if len(scopes) == 0 {
		return resources
	}
	set := make(map[string]struct{}, len(scopes))
	for _, s := range scopes {
		key := s
		if !opts.CaseSensitive {
			key = strings.ToLower(s)
		}
		set[key] = struct{}{}
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["scope"]
		if !ok {
			continue
		}
		if !opts.CaseSensitive {
			v = strings.ToLower(v)
		}
		if _, found := set[v]; found {
			out = append(out, r)
		}
	}
	return out
}
