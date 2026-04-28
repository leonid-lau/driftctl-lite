package tfstate

import (
	"strings"

	"github.com/snyk/driftctl-lite/internal/tfstate/resource"
)

// EndpointFilterOptions configures endpoint filtering behaviour.
type EndpointFilterOptions struct {
	CaseSensitive bool
}

// DefaultEndpointFilterOptions returns sensible defaults (case-insensitive).
func DefaultEndpointFilterOptions() EndpointFilterOptions {
	return EndpointFilterOptions{CaseSensitive: false}
}

// FilterByEndpoint returns resources whose "endpoint" attribute matches the
// given value. An empty endpoint string returns all resources unchanged.
func FilterByEndpoint(resources []resource.Resource, endpoint string, opts EndpointFilterOptions) []resource.Resource {
	if endpoint == "" {
		return resources
	}
	var out []resource.Resource
	for _, r := range resources {
		val, ok := r.Attributes["endpoint"]
		if !ok {
			continue
		}
		if matchEndpoint(val, endpoint, opts.CaseSensitive) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByEndpoints returns resources matching ANY of the provided endpoints
// (OR semantics).
func FilterByEndpoints(resources []resource.Resource, endpoints []string, opts EndpointFilterOptions) []resource.Resource {
	if len(endpoints) == 0 {
		return resources
	}
	var out []resource.Resource
	for _, r := range resources {
		val, ok := r.Attributes["endpoint"]
		if !ok {
			continue
		}
		for _, ep := range endpoints {
			if matchEndpoint(val, ep, opts.CaseSensitive) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchEndpoint(val, target string, caseSensitive bool) bool {
	if caseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
