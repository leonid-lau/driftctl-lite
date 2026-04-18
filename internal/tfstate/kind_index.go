package tfstate

import "strings"

// KindIndex maps normalised kind strings to slices of resources.
type KindIndex struct {
	index map[string][]Resource
}

// BuildKindIndex constructs a KindIndex from the given resources.
// Keys are stored in lower-case for case-insensitive lookup.
func BuildKindIndex(resources []Resource) *KindIndex {
	idx := &KindIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		key := strings.ToLower(r.Kind)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns all resources with the given kind (case-insensitive).
// Returns nil when no resources match.
func (ki *KindIndex) Lookup(kind string) []Resource {
	if kind == "" {
		return nil
	}
	return ki.index[strings.ToLower(kind)]
}

// Kinds returns all distinct kind strings present in the index.
func (ki *KindIndex) Kinds() []string {
	out := make([]string, 0, len(ki.index))
	for k := range ki.index {
		out = append(out, k)
	}
	return out
}
