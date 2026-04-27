package tfstate

import "strings"

// PortIndex maps normalised port values to matching resources.
type PortIndex struct {
	index map[string][]Resource
}

// BuildPortIndex constructs a case-insensitive index over the "port" attribute.
func BuildPortIndex(resources []Resource) *PortIndex {
	idx := &PortIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["port"]
		if !ok {
			continue
		}
		key := strings.ToLower(toString(v))
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns all resources whose port matches (case-insensitive).
func (idx *PortIndex) Lookup(port string) []Resource {
	return idx.index[strings.ToLower(port)]
}

// Ports returns all distinct (normalised) port values present in the index.
func (idx *PortIndex) Ports() []string {
	keys := make([]string, 0, len(idx.index))
	for k := range idx.index {
		keys = append(keys, k)
	}
	return keys
}
