package tfstate

import "strings"

// CategoryIndex maps lowercase category values to slices of resources.
type CategoryIndex struct {
	index map[string][]Resource
}

// BuildCategoryIndex constructs a CategoryIndex from the given resources.
// Keys are stored in lowercase to support case-insensitive lookup.
func BuildCategoryIndex(resources []Resource) *CategoryIndex {
	idx := &CategoryIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["category"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok || s == "" {
			continue
		}
		key := strings.ToLower(s)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns all resources for the given category (case-insensitive).
// Returns nil if no resources are found.
func (idx *CategoryIndex) Lookup(category string) []Resource {
	if category == "" {
		return nil
	}
	return idx.index[strings.ToLower(category)]
}

// Categories returns all distinct category keys present in the index.
func (idx *CategoryIndex) Categories() []string {
	keys := make([]string, 0, len(idx.index))
	for k := range idx.index {
		keys = append(keys, k)
	}
	return keys
}
