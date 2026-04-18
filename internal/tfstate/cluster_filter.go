package tfstate

import "strings"

// DefaultClusterFilterOptions returns case-insensitive matching by default.
type ClusterFilterOptions struct {
	CaseSensitive bool
}

func DefaultClusterFilterOptions() ClusterFilterOptions {
	return ClusterFilterOptions{CaseSensitive: false}
}

// FilterByCluster returns resources whose Cluster metadata field matches cluster.
// An empty cluster string returns all resources unchanged.
func FilterByCluster(resources []Resource, cluster string, opts ClusterFilterOptions) []Resource {
	if cluster == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Metadata["cluster"]
		if !ok {
			continue
		}
		if matchCluster(v, cluster, opts.CaseSensitive) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByClusters returns resources matching ANY of the provided cluster names (OR semantics).
func FilterByClusters(resources []Resource, clusters []string, opts ClusterFilterOptions) []Resource {
	if len(clusters) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Metadata["cluster"]
		if !ok {
			continue
		}
		for _, c := range clusters {
			if matchCluster(v, c, opts.CaseSensitive) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchCluster(val, target string, caseSensitive bool) bool {
	if caseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
