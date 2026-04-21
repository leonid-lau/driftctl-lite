package tfstate

import "strings"

// DefaultPolicyFilterOptions returns a PolicyFilterOptions with CaseSensitive=false.
func DefaultPolicyFilterOptions() PolicyFilterOptions {
	return PolicyFilterOptions{CaseSensitive: false}
}

// PolicyFilterOptions controls matching behaviour for policy filters.
type PolicyFilterOptions struct {
	CaseSensitive bool
}

// FilterByPolicy returns resources whose "policy" attribute matches the given value.
// An empty policy string returns all resources unchanged.
func FilterByPolicy(resources []Resource, policy string, opts PolicyFilterOptions) []Resource {
	if policy == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["policy"]
		if !ok {
			continue
		}
		if matchPolicy(v, policy, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByPolicies returns resources matching ANY of the given policy values (OR semantics).
func FilterByPolicies(resources []Resource, policies []string, opts PolicyFilterOptions) []Resource {
	if len(policies) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["policy"]
		if !ok {
			continue
		}
		for _, p := range policies {
			if matchPolicy(v, p, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchPolicy(attrVal, target string, opts PolicyFilterOptions) bool {
	if opts.CaseSensitive {
		return attrVal == target
	}
	return strings.EqualFold(attrVal, target)
}
