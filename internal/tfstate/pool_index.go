package tfstate

import "strings"

// PoolIndex maps normalised pool names to slices of resources.
type PoolIndex struct {
	index map[string][]Resource
}

// BuildPoolIndex constructs a PoolIndex from the given resources.
// Keys are stored in lower-case for case-insensitive lookup.
func BuildPoolIndex(resources []Resource) *PoolIndex {
	idx := &PoolIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		if v, ok := r.Attributes["pool"]; ok && v != "" {
			key := strings.ToLower(v)
			idx.index[key] = append(idx.index[key], r)
		}
	}
	return idx
}

// Lookup returns all resources that belong to the given pool (case-insensitive).
func (i *PoolIndex) Lookup(pool string) []Resource {
	if pool == "" {
		return nil
	}
	return i.index[strings.ToLower(pool)]
}

// Pools returns all distinct pool names stored in the index.
func (i *PoolIndex) Pools() []string {
	keys := make([]string, 0, len(i.index))
	for k := range i.index {
		keys = append(keys, k)
	}
	return keys
}
