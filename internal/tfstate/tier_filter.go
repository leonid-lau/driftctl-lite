package tfstate

import "strings"

// TierFilterOptions configures FilterByTier behaviour.
type TierFilterOptions struct {
	CaseSensitive bool
}

// DefaultTierFilterOptions returns sensible defaults (case-insensitive).
func DefaultTierFilterOptions() TierFilterOptions {
	return TierFilterOptions{CaseSensitive: false}
}

// FilterByTier returns resources whose "tier" attribute matches the given value.
// An empty tier string returns all resources unchanged.
func FilterByTier(resources []Resource, tier string, opts TierFilterOptions) []Resource {
	if tier == "" {
		return resources
	}
	want := tier
	if !opts.CaseSensitive {
		want = strings.ToLower(tier)
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["tier"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		candidate := s
		if !opts.CaseSensitive {
			candidate = strings.ToLower(s)
		}
		if candidate == want {
			out = append(out, r)
		}
	}
	return out
}

// FilterByTiers returns resources matching ANY of the provided tiers (OR semantics).
func FilterByTiers(resources []Resource, tiers []string, opts TierFilterOptions) []Resource {
	if len(tiers) == 0 {
		return resources
	}
	var out []Resource
	for _, t := range tiers {
		out = append(out, FilterByTier(resources, t, opts)...)
	}
	return out
}
