package tfstate

import "strings"

// DefaultCostFilterOptions returns case-insensitive options.
type CostFilterOptions struct {
	CaseSensitive bool
}

func DefaultCostFilterOptions() CostFilterOptions {
	return CostFilterOptions{CaseSensitive: false}
}

// FilterByCost returns resources whose "cost_center" attribute matches the given value.
func FilterByCost(resources []Resource, costCenter string, opts CostFilterOptions) []Resource {
	if costCenter == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["cost_center"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		if matchCost(s, costCenter, opts.CaseSensitive) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByCosts returns resources matching any of the provided cost centers (OR semantics).
func FilterByCosts(resources []Resource, costCenters []string, opts CostFilterOptions) []Resource {
	if len(costCenters) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["cost_center"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		for _, cc := range costCenters {
			if matchCost(s, cc, opts.CaseSensitive) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchCost(val, target string, caseSensitive bool) bool {
	if !caseSensitive {
		return strings.EqualFold(val, target)
	}
	return val == target
}
