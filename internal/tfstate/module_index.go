package tfstate

import "strings"

// ModuleIndex maps lowercase module name -> list of resources.
type ModuleIndex struct {
	index map[string][]Resource
}

// BuildModuleIndex constructs a ModuleIndex from a slice of resources.
func BuildModuleIndex(resources []Resource) *ModuleIndex {
	idx := &ModuleIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, _ := r.Attributes["module"].(string)
		key := strings.ToLower(v)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources for the given module (case-insensitive).
func (i *ModuleIndex) Lookup(module string) []Resource {
	return i.index[strings.ToLower(module)]
}

// Modules returns all known module names in the index.
func (i *ModuleIndex) Modules() []string {
	keys := make([]string, 0, len(i.index))
	for k := range i.index {
		keys = append(keys, k)
	}
	return keys
}
