package tfstate

import "strings"

// EnvironmentIndex maps environment name -> list of resources.
type EnvironmentIndex struct {
	index map[string][]Resource
}

// BuildEnvironmentIndex builds an index from a slice of resources.
// Keys are lowercased for case-insensitive lookup.
func BuildEnvironmentIndex(resources []Resource) *EnvironmentIndex {
	idx := &EnvironmentIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		val, ok := r.Attributes["environment"]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		key := strings.ToLower(s)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources for the given environment (case-insensitive).
func (e *EnvironmentIndex) Lookup(env string) []Resource {
	return e.index[strings.ToLower(env)]
}

// Environments returns all indexed environment names.
func (e *EnvironmentIndex) Environments() []string {
	keys := make([]string, 0, len(e.index))
	for k := range e.index {
		keys = append(keys, k)
	}
	return keys
}
