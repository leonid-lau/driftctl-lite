package tfstate

import "strings"

// ComponentIndex maps lowercase component name to resources.
type ComponentIndex struct {
	index map[string][]Resource
}

// BuildComponentIndex builds an index from a slice of resources.
func BuildComponentIndex(resources []Resource) *ComponentIndex {
	idx := &ComponentIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["component"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		key := strings.ToLower(s)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources for the given component (case-insensitive).
func (i *ComponentIndex) Lookup(component string) []Resource {
	return i.index[strings.ToLower(component)]
}

// Components returns all indexed component names.
func (i *ComponentIndex) Components() []string {
	keys := make([]string, 0, len(i.index))
	for k := range i.index {
		keys = append(keys, k)
	}
	return keys
}
