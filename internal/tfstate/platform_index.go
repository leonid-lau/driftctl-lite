package tfstate

import "strings"

// PlatformIndex maps lower-cased platform names to the resources that carry them.
type PlatformIndex struct {
	index map[string][]Resource
}

// BuildPlatformIndex constructs a PlatformIndex from the given resource slice.
func BuildPlatformIndex(resources []Resource) *PlatformIndex {
	idx := &PlatformIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["platform"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok || s == "" {
			continue
		}
		key := strings.ToLower(s)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources whose platform matches the given value
// (case-insensitive). Returns nil when no match is found.
func (pi *PlatformIndex) Lookup(platform string) []Resource {
	if platform == "" {
		return nil
	}
	return pi.index[strings.ToLower(platform)]
}

// Platforms returns all distinct platform names present in the index.
func (pi *PlatformIndex) Platforms() []string {
	keys := make([]string, 0, len(pi.index))
	for k := range pi.index {
		keys = append(keys, k)
	}
	return keys
}
