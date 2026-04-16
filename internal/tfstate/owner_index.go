package tfstate

// OwnerIndex maps owner name -> list of resources.
type OwnerIndex struct {
	index map[string][]Resource
}

// BuildOwnerIndex constructs an index from a slice of resources.
func BuildOwnerIndex(resources []Resource) *OwnerIndex {
	idx := &OwnerIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		owner := toLower(r.Metadata["owner"])
		idx.index[owner] = append(idx.index[owner], r)
	}
	return idx
}

// Lookup returns all resources for the given owner (case-insensitive).
func (i *OwnerIndex) Lookup(owner string) []Resource {
	if owner == "" {
		return nil
	}
	return i.index[toLower(owner)]
}

// Owners returns all distinct owner values present in the index.
func (i *OwnerIndex) Owners() []string {
	owners := make([]string, 0, len(i.index))
	for k := range i.index {
		owners = append(owners, k)
	}
	return owners
}
