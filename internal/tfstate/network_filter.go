package tfstate

import "strings"

// DefaultNetworkFilterOptions returns case-insensitive matching by default.
type NetworkFilterOptions struct {
	CaseSensitive bool
}

func DefaultNetworkFilterOptions() NetworkFilterOptions {
	return NetworkFilterOptions{CaseSensitive: false}
}

// FilterByNetwork returns resources whose "network" attribute matches the given value.
func FilterByNetwork(resources []Resource, network string, opts NetworkFilterOptions) []Resource {
	if network == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchNetwork(r, network, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByNetworks returns resources matching ANY of the provided network values (OR semantics).
func FilterByNetworks(resources []Resource, networks []string, opts NetworkFilterOptions) []Resource {
	if len(networks) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, n := range networks {
			if matchNetwork(r, n, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchNetwork(r Resource, network string, opts NetworkFilterOptions) bool {
	v, ok := r.Attributes["network"]
	if !ok {
		v, ok = r.Attributes["vpc"]
	}
	if !ok {
		return false
	}
	s, _ := v.(string)
	if opts.CaseSensitive {
		return s == network
	}
	return strings.EqualFold(s, network)
}
