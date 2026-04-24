package tfstate

import "strings"

// ArchFilterOptions controls matching behaviour for architecture filters.
type ArchFilterOptions struct {
	CaseSensitive bool
}

// DefaultArchFilterOptions returns case-insensitive options.
func DefaultArchFilterOptions() ArchFilterOptions {
	return ArchFilterOptions{CaseSensitive: false}
}

// FilterByArch returns resources whose "arch" or "architecture" attribute
// matches the given value. An empty arch returns all resources unchanged.
func FilterByArch(resources []Resource, arch string, opts ArchFilterOptions) []Resource {
	if arch == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchArch(r, arch, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByArchs returns resources matching ANY of the provided architectures
// (OR semantics).
func FilterByArchs(resources []Resource, archs []string, opts ArchFilterOptions) []Resource {
	if len(archs) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, a := range archs {
			if matchArch(r, a, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchArch(r Resource, arch string, opts ArchFilterOptions) bool {
	keys := []string{"arch", "architecture"}
	for _, k := range keys {
		if v, ok := r.Attributes[k]; ok {
			s, ok2 := v.(string)
			if !ok2 {
				continue
			}
			if opts.CaseSensitive {
				if s == arch {
					return true
				}
			} else {
				if strings.EqualFold(s, arch) {
					return true
				}
			}
		}
	}
	return false
}
