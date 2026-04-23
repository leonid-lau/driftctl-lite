package tfstate

import "strings"

// SourceIndex maps lowercase source name -> list of resources.
type SourceIndex struct {
	index map[string][]Resource
}

// BuildSourceIndex constructs a SourceIndex from a slice of resources.
func BuildSourceIndex(resources []Resource) *SourceIndex {
	idx := &SourceIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		key := strings.ToLower(r.Source)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns all resources with the given source (case-insensitive).
func (si *SourceIndex) Lookup(source string) []Resource {
	return si.index[strings.ToLower(source)]
}

// Sources returns all distinct source values present in the index.
func (si *SourceIndex) Sources() []string {
	keys := make([]string, 0, len(si.index))
	for k := range si.index {
		keys = append(keys, k)
	}
	return keys
}

// Count returns the total number of resources across all sources.
func (si *SourceIndex) Count() int {
	total := 0
	for _, resources := range si.index {
		total += len(resources)
	}
	return total
}
