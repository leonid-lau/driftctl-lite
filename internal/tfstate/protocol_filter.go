package tfstate

import "strings"

// DefaultProtocolFilterOptions returns case-insensitive matching options.
func DefaultProtocolFilterOptions() FilterOptions {
	return FilterOptions{CaseSensitive: false}
}

// FilterOptions controls matching behaviour for protocol filters.
type FilterOptions struct {
	CaseSensitive bool
}

// FilterByProtocol returns resources whose "protocol" attribute matches the
// given value. An empty protocol returns all resources unchanged.
func FilterByProtocol(resources []Resource, protocol string, opts FilterOptions) []Resource {
	if protocol == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchProtocol(r, protocol, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByProtocols returns resources matching ANY of the given protocols
// (OR semantics). An empty slice returns all resources unchanged.
func FilterByProtocols(resources []Resource, protocols []string, opts FilterOptions) []Resource {
	if len(protocols) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, p := range protocols {
			if matchProtocol(r, p, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchProtocol(r Resource, protocol string, opts FilterOptions) bool {
	v, ok := r.Attributes["protocol"]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	if opts.CaseSensitive {
		return s == protocol
	}
	return strings.EqualFold(s, protocol)
}
