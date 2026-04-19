package tfstate

import "strings"

// DefaultServiceFilterOptions returns case-insensitive matching options.
type ServiceFilterOptions struct {
	CaseSensitive bool
}

func DefaultServiceFilterOptions() ServiceFilterOptions {
	return ServiceFilterOptions{CaseSensitive: false}
}

// FilterByService returns resources whose "service" attribute matches the given value.
// If service is empty, all resources are returned.
func FilterByService(resources []Resource, service string, opts ServiceFilterOptions) []Resource {
	if service == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["service"]
		if !ok {
			continue
		}
		s, _ := v.(string)
		if matchService(s, service, opts.CaseSensitive) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByServices returns resources matching any of the given service values (OR semantics).
func FilterByServices(resources []Resource, services []string, opts ServiceFilterOptions) []Resource {
	if len(services) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["service"]
		if !ok {
			continue
		}
		s, _ := v.(string)
		for _, svc := range services {
			if matchService(s, svc, opts.CaseSensitive) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchService(val, target string, caseSensitive bool) bool {
	if caseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
