package tfstate

// NamespaceIndex maps namespace -> list of resources.
type NamespaceIndex struct {
	index map[string][]Resource
}

// BuildNamespaceIndex constructs an index from a slice of resources.
func BuildNamespaceIndex(resources []Resource) *NamespaceIndex {
	idx := &NamespaceIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		ns := r.Namespace
		idx.index[ns] = append(idx.index[ns], r)
	}
	return idx
}

// Lookup returns all resources in the given namespace.
func (idx *NamespaceIndex) Lookup(ns string) []Resource {
	return idx.index[ns]
}

// Namespaces returns all distinct namespaces present in the index.
func (idx *NamespaceIndex) Namespaces() []string {
	keys := make([]string, 0, len(idx.index))
	for k := range idx.index {
		keys = append(keys, k)
	}
	return keys
}
