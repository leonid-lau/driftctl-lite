package tfstate

import (
	"strings"

	"github.com/snyk/driftctl-lite/internal/tfstate"
)

// DefaultCIDRFilterOptions returns case-insensitive matching options.
type CIDRFilterOptions struct {
	CaseSensitive bool
}

func DefaultCIDRFilterOptions() CIDRFilterOptions {
	return CIDRFilterOptions{CaseSensitive: false}
}

// FilterByCIDR returns resources whose "cidr" or "cidr_block" attribute matches
// the given CIDR string. An empty cidr returns all resources unchanged.
func FilterByCIDR(resources []Resource, cidr string, opts CIDRFilterOptions) []Resource {
	if cidr == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchCIDR(r, cidr, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByCIDRs returns resources that match ANY of the provided CIDRs (OR semantics).
func FilterByCIDRs(resources []Resource, cidrs []string, opts CIDRFilterOptions) []Resource {
	if len(cidrs) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, c := range cidrs {
			if matchCIDR(r, c, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchCIDR(r Resource, cidr string, opts CIDRFilterOptions) bool {
	candidates := []string{
		stringAttr(r, "cidr"),
		stringAttr(r, "cidr_block"),
	}
	for _, v := range candidates {
		if v == "" {
			continue
		}
		if opts.CaseSensitive {
			if v == cidr {
				return true
			}
		} else {
			if strings.EqualFold(v, cidr) {
				return true
			}
		}
	}
	return false
}

func stringAttr(r Resource, key string) string {
	if r.Attributes == nil {
		return ""
	}
	v, _ := r.Attributes[key].(string)
	return v
}
