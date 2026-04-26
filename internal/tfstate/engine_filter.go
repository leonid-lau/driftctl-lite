package tfstate

import "strings"

// EngineFilterOptions controls matching behaviour for engine filters.
type EngineFilterOptions struct {
	// CaseSensitive disables case-folding when true.
	CaseSensitive bool
}

// DefaultEngineFilterOptions returns the recommended defaults.
func DefaultEngineFilterOptions() EngineFilterOptions {
	return EngineFilterOptions{CaseSensitive: false}
}

// FilterByEngine returns resources whose "engine" attribute equals the given
// value. An empty engine string returns the full slice unchanged.
func FilterByEngine(resources []Resource, engine string, opts EngineFilterOptions) []Resource {
	if engine == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchEngine(r, engine, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByEngines returns resources that match ANY of the supplied engines
// (OR semantics).
func FilterByEngines(resources []Resource, engines []string, opts EngineFilterOptions) []Resource {
	if len(engines) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, e := range engines {
			if matchEngine(r, e, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchEngine(r Resource, engine string, opts EngineFilterOptions) bool {
	v, ok := r.Attributes["engine"]
	if !ok {
		return false
	}
	got, _ := v.(string)
	if opts.CaseSensitive {
		return got == engine
	}
	return strings.EqualFold(got, engine)
}
