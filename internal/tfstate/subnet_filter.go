package tfstate

import "strings"

// SubnetFilterOptions controls matching behaviour for FilterBySubnet.
type SubnetFilterOptions struct {
	CaseSensitive bool
}

// DefaultSubnetFilterOptions returns sensible defaults (case-insensitive).
func DefaultSubnetFilterOptions() SubnetFilterOptions {
	return SubnetFilterOptions{CaseSensitive: false}
}

// FilterBySubnet returns resources whose "subnet" attribute matches the given
// value. An empty subnet string returns all resources unchanged.
func FilterBySubnet(resources []Resource, subnet string, opts SubnetFilterOptions) []Resource {
	if subnet == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchSubnet(r, subnet, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterBySubnets returns resources that match ANY of the provided subnet
// values (OR semantics).
func FilterBySubnets(resources []Resource, subnets []string, opts SubnetFilterOptions) []Resource {
	if len(subnets) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, s := range subnets {
			if matchSubnet(r, s, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchSubnet(r Resource, subnet string, opts SubnetFilterOptions) bool {
	v, ok := r.Attributes["subnet"]
	if !ok {
		v, ok = r.Attributes["subnet_id"]
	}
	if !ok {
		return false
	}
	if opts.CaseSensitive {
		return v == subnet
	}
	return strings.EqualFold(v, subnet)
}
