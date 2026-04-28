package tfstate

import "strings"

// ArchIndex maps architecture values to matching resources.
type ArchIndex struct {
	index map[string][]Resource
}

// BuildArchIndex constructs an ArchIndex from a slice of resources.
// Lookup keys are normalised to lower-case.
func BuildArchIndex(resources []Resource) *ArchIndex {
	idx := &ArchIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		arch := strings.TrimSpace(r.Attributes["arch"])
		if arch == "" {
			arch = strings.TrimSpace(r.Attributes["architecture"])
		}
		if arch == "" {
			continue
		}
		key := strings.ToLower(arch)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources matching the given architecture (case-insensitive).
// Returns nil when the key is empty or not found.
func (idx *ArchIndex) Lookup(arch string) []Resource {
	if arch == "" {
		return nil
	}
	return idx.index[strings.ToLower(arch)]
}

// Archs returns all distinct architecture keys stored in the index.
func (idx *ArchIndex) Archs() []string {
	keys := make([]string, 0, len(idx.index))
	for k := range idx.index {
		keys = append(keys, k)
	}
	return keys
}
