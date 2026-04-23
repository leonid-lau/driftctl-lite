package tfstate

import "strings"

// DefaultPlatformFilterOptions returns a PlatformFilterOptions with CaseSensitive=false.
func DefaultPlatformFilterOptions() PlatformFilterOptions {
	return PlatformFilterOptions{CaseSensitive: false}
}

// PlatformFilterOptions controls matching behaviour for platform filters.
type PlatformFilterOptions struct {
	CaseSensitive bool
}

// FilterByPlatform returns resources whose "platform" attribute matches the
// given value. An empty platform string returns all resources unchanged.
func FilterByPlatform(resources []Resource, platform string, opts PlatformFilterOptions) []Resource {
	if platform == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchPlatform(r, platform, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByPlatforms returns resources that match ANY of the supplied platform
// values (OR semantics). An empty slice returns all resources unchanged.
func FilterByPlatforms(resources []Resource, platforms []string, opts PlatformFilterOptions) []Resource {
	if len(platforms) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, p := range platforms {
			if matchPlatform(r, p, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchPlatform(r Resource, platform string, opts PlatformFilterOptions) bool {
	v, ok := r.Attributes["platform"]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	if opts.CaseSensitive {
		return s == platform
	}
	return strings.EqualFold(s, platform)
}
