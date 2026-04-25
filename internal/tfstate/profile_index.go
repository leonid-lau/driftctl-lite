package tfstate

import "strings"

// ProfileIndex maps lowercase profile values to slices of Resources.
type ProfileIndex struct {
	index map[string][]Resource
}

// BuildProfileIndex constructs a ProfileIndex from the given resources.
func BuildProfileIndex(resources []Resource) *ProfileIndex {
	idx := &ProfileIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["profile"]
		if !ok || v == "" {
			continue
		}
		key := strings.ToLower(v)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources matching the given profile (case-insensitive).
func (pi *ProfileIndex) Lookup(profile string) []Resource {
	if profile == "" {
		return nil
	}
	return pi.index[strings.ToLower(profile)]
}

// Profiles returns all distinct profile values present in the index.
func (pi *ProfileIndex) Profiles() []string {
	keys := make([]string, 0, len(pi.index))
	for k := range pi.index {
		keys = append(keys, k)
	}
	return keys
}
