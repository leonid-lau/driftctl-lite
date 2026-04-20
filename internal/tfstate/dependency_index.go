package tfstate

// DependencyIndex maps dependency name -> list of resources.
type DependencyIndex struct {
	index map[string][]Resource
}

// BuildDependencyIndex constructs a DependencyIndex from a slice of resources.
// Keys are stored in lower-case for case-insensitive look-up.
func BuildDependencyIndex(resources []Resource) *DependencyIndex {
	idx := &DependencyIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["dependency"]
		if !ok {
			continue
		}
		s, _ := v.(string)
		if s == "" {
			continue
		}
		key := toLower(s)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources with the given dependency (case-insensitive).
func (i *DependencyIndex) Lookup(dep string) []Resource {
	return i.index[toLower(dep)]
}

// Dependencies returns all unique dependency names stored in the index.
func (i *DependencyIndex) Dependencies() []string {
	keys := make([]string, 0, len(i.index))
	for k := range i.index {
		keys = append(keys, k)
	}
	return keys
}
