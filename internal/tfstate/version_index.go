package tfstate

// VersionIndex maps version string to resources.
type VersionIndex struct {
	index map[string][]Resource
}

// BuildVersionIndex constructs a VersionIndex from a slice of resources.
func BuildVersionIndex(resources []Resource) *VersionIndex {
	idx := &VersionIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		key := toLower(r.Version)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources matching the given version (case-insensitive).
func (idx *VersionIndex) Lookup(version string) []Resource {
	if version == "" {
		return nil
	}
	return idx.index[toLower(version)]
}

// Versions returns all indexed version strings.
func (idx *VersionIndex) Versions() []string {
	keys := make([]string, 0, len(idx.index))
	for k := range idx.index {
		keys = append(keys, k)
	}
	return keys
}
