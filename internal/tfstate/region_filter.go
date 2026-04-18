package tfstate

import "strings"

// DefaultRegionFilterOptions returns options with case-insensitive matching enabled.
func DefaultRegionFilterOptions() RegionFilterOptions {
	return RegionFilterOptions{CaseInsensitive: true}
}

// RegionFilterOptions controls behaviour of region filtering.
type RegionFilterOptions struct {
	CaseInsensitive bool
}

// FilterByRegion returns resources whose Region field matches the given region.
// An empty region string returns all resources unchanged.
func FilterByRegion(resources []Resource, region string, opts RegionFilterOptions) []Resource {
	if region == "" {
		return resources
	}
	needle := region
	if opts.CaseInsensitive {
		needle = strings.ToLower(region)
	}
	var out []Resource
	for _, r := range resources {
		haystack := r.Region
		if opts.CaseInsensitive {
			haystack = strings.ToLower(r.Region)
		}
		if haystack == needle {
			out = append(out, r)
		}
	}
	return out
}

// FilterByRegions returns resources matching ANY of the provided regions (OR semantics).
// An empty slice returns all resources unchanged.
func FilterByRegions(resources []Resource, regions []string, opts RegionFilterOptions) []Resource {
	if len(regions) == 0 {
		return resources
	}
	needles := make(map[string]struct{}, len(regions))
	for _, r := range regions {
		k := r
		if opts.CaseInsensitive {
			k = strings.ToLower(r)
		}
		needles[k] = struct{}{}
	}
	var out []Resource
	for _, r := range resources {
		haystack := r.Region
		if opts.CaseInsensitive {
			haystack = strings.ToLower(r.Region)
		}
		if _, ok := needles[haystack]; ok {
			out = append(out, r)
		}
	}
	return out
}
