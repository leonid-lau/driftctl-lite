package tfstate

import "strings"

// DefaultPortFilterOptions returns a PortFilterOptions with CaseSensitive=false.
func DefaultPortFilterOptions() PortFilterOptions {
	return PortFilterOptions{CaseSensitive: false}
}

// PortFilterOptions controls how port matching is performed.
type PortFilterOptions struct {
	CaseSensitive bool
}

// FilterByPort returns resources whose "port" attribute matches the given value.
// An empty port string returns all resources unchanged.
func FilterByPort(resources []Resource, port string, opts PortFilterOptions) []Resource {
	if port == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchPort(r, port, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByPorts returns resources matching ANY of the given port values (OR semantics).
func FilterByPorts(resources []Resource, ports []string, opts PortFilterOptions) []Resource {
	if len(ports) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, p := range ports {
			if matchPort(r, p, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchPort(r Resource, port string, opts PortFilterOptions) bool {
	v, ok := r.Attributes["port"]
	if !ok {
		return false
	}
	val := toString(v)
	if opts.CaseSensitive {
		return val == port
	}
	return strings.EqualFold(val, port)
}
