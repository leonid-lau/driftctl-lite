package tfstate

import (
	"strings"

	"github.com/snyk/driftctl-lite/internal/tfstate"
)

// DefaultVersionFilterOptions returns case-insensitive options.
type VersionFilterOptions struct {
	CaseSensitive bool
}

func DefaultVersionFilterOptions() VersionFilterOptions {
	return VersionFilterOptions{CaseSensitive: false}
}

// FilterByVersion returns resources whose Version field matches version.
// An empty version string returns all resources unchanged.
func FilterByVersion(resources []tfstate.Resource, version string, opts VersionFilterOptions) []tfstate.Resource {
	if version == "" {
		return resources
	}
	var out []tfstate.Resource
	for _, r := range resources {
		v := r.Attributes["version"]
		if match := func(a, b string) bool {
			if opts.CaseSensitive {
				return a == b
			}
			return strings.EqualFold(a, b)
		}(v, version); match {
			out = append(out, r)
		}
	}
	return out
}

// FilterByVersions returns resources matching any of the provided versions (OR semantics).
func FilterByVersions(resources []tfstate.Resource, versions []string, opts VersionFilterOptions) []tfstate.Resource {
	if len(versions) == 0 {
		return resources
	}
	var out []tfstate.Resource
	for _, r := range resources {
		v := r.Attributes["version"]
		for _, ver := range versions {
			if opts.CaseSensitive {
				if v == ver {
					out = append(out, r)
					break
				}
			} else {
				if strings.EqualFold(v, ver) {
					out = append(out, r)
					break
				}
			}
		}
	}
	return out
}
