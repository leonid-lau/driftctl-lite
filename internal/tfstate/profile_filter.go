package tfstate

import "strings"

// DefaultProfileFilterOptions returns a FilterOptions with case-insensitive matching.
type ProfileFilterOptions struct {
	CaseSensitive bool
}

// DefaultProfileFilterOptions returns sensible defaults.
func DefaultProfileFilterOptions() ProfileFilterOptions {
	return ProfileFilterOptions{CaseSensitive: false}
}

// FilterByProfile returns resources whose "profile" attribute matches the given value.
// If profile is empty, all resources are returned.
func FilterByProfile(resources []Resource, profile string, opts ProfileFilterOptions) []Resource {
	if profile == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["profile"]
		if !ok {
			continue
		}
		if matchProfile(v, profile, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByProfiles returns resources matching ANY of the given profile values (OR semantics).
func FilterByProfiles(resources []Resource, profiles []string, opts ProfileFilterOptions) []Resource {
	if len(profiles) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["profile"]
		if !ok {
			continue
		}
		for _, p := range profiles {
			if matchProfile(v, p, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchProfile(val, target string, opts ProfileFilterOptions) bool {
	if opts.CaseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
