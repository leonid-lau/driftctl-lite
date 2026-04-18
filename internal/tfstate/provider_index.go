package tfstate

// ProviderIndex maps provider name to resources.
type ProviderIndex struct {
	index map[string][]Resource
}

// BuildProviderIndex constructs a ProviderIndex from a slice of resources.
func BuildProviderIndex(resources []Resource) *ProviderIndex {
	idx := &ProviderIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		key := toLower(r.Provider)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources matching the given provider (case-insensitive).
func (idx *ProviderIndex) Lookup(provider string) []Resource {
	if provider == "" {
		return nil
	}
	return idx.index[toLower(provider)]
}

// Providers returns all indexed provider names.
func (idx *ProviderIndex) Providers() []string {
	keys := make([]string, 0, len(idx.index))
	for k := range idx.index {
		keys = append(keys, k)
	}
	return keys
}
