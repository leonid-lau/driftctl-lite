package tfstate

// ClusterIndex maps cluster name -> list of resources.
type ClusterIndex struct {
	index map[string][]Resource
}

// BuildClusterIndex constructs a ClusterIndex from a slice of resources.
// Keys are stored in their original case; lookups are caller's responsibility.
func BuildClusterIndex(resources []Resource) *ClusterIndex {
	idx := &ClusterIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		if v, ok := r.Metadata["cluster"]; ok && v != "" {
			idx.index[v] = append(idx.index[v], r)
		}
	}
	return idx
}

// Lookup returns resources for the given cluster name (exact match).
func (ci *ClusterIndex) Lookup(cluster string) []Resource {
	return ci.index[cluster]
}

// Clusters returns all indexed cluster names.
func (ci *ClusterIndex) Clusters() []string {
	keys := make([]string, 0, len(ci.index))
	for k := range ci.index {
		keys = append(keys, k)
	}
	return keys
}
