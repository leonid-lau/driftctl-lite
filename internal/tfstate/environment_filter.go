package tfstate

import "strings"

// DefaultEnvironmentFilterOptions returns case-insensitive matching by default.
type EnvironmentFilterOptions struct {
	CaseSensitive bool
}

func DefaultEnvironmentFilterOptions() EnvironmentFilterOptions {
	return EnvironmentFilterOptions{CaseSensitive: false}
}

// FilterByEnvironment returns resources whose "environment" attribute matches env.
// If env is empty, all resources are returned.
func FilterByEnvironment(resources []Resource, env string, opts EnvironmentFilterOptions) []Resource {
	if env == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		val, ok := r.Attributes["environment"]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		if opts.CaseSensitive {
			if s == env {
				out = append(out, r)
			}
		} else {
			if strings.EqualFold(s, env) {
				out = append(out, r)
			}
		}
	}
	return out
}

// FilterByEnvironments returns resources matching any of the given environments (OR semantics).
func FilterByEnvironments(resources []Resource, envs []string, opts EnvironmentFilterOptions) []Resource {
	if len(envs) == 0 {
		return resources
	}
	set := make(map[string]struct{}, len(envs))
	for _, e := range envs {
		key := e
		if !opts.CaseSensitive {
			key = strings.ToLower(e)
		}
		set[key] = struct{}{}
	}
	var out []Resource
	for _, r := range resources {
		val, ok := r.Attributes["environment"]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		key := s
		if !opts.CaseSensitive {
			key = strings.ToLower(s)
		}
		if _, found := set[key]; found {
			out = append(out, r)
		}
	}
	return out
}
