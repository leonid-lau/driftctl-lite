package tfstate

import "strings"

// RegionIndex maps region name (lowercase) -> list of resources.
type RegionIndex struct {
	index map[string][]Resource
}

// BuildRegionIndex constructs a RegionIndex from a slice of resources.
// Keys are stored in lowercase for case-insensitive lookups.
func BuildRegionIndex(resources []Resource) *RegionIndex {
	idx := &RegionIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		if v, ok := r.Metadata["region"]; ok && v != "" {
			key := strings.ToLower(v)
			idx.index[key] = append(idx.index[key], r)
		}
	}
	return idx
}

// Lookup returns resources for the given region (case-insensitive).
func (ri *RegionIndex) Lookup(region string) []Resource {
	return ri.index[strings.ToLower(region)]
}

// Regions returns all indexed region names (lowercase).
func (ri *RegionIndex) Regions() []string {
	keys := make([]string, 0, len(ri.index))
	for k := range ri.index {
		keys = append(keys, k)
	}
	return keys
}
