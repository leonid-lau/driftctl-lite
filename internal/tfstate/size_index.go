package tfstate

import "strings"

// SizeIndex maps a normalised size value to the resources that carry it.
type SizeIndex struct {
	m map[string][]Resource
}

// BuildSizeIndex constructs a case-insensitive index over the "size" attribute
// of the supplied resources.
func BuildSizeIndex(resources []Resource) *SizeIndex {
	idx := &SizeIndex{m: make(map[string][]Resource)}
	for _, r := range resources {
		v, _ := r.Attributes["size"].(string)
		if v == "" {
			continue
		}
		key := strings.ToLower(v)
		idx.m[key] = append(idx.m[key], r)
	}
	return idx
}

// Lookup returns resources whose size matches the given value
// (case-insensitive). Returns nil when no match exists.
func (idx *SizeIndex) Lookup(size string) []Resource {
	if size == "" {
		return nil
	}
	return idx.m[strings.ToLower(size)]
}

// Sizes returns the distinct (normalised) size values present in the index.
func (idx *SizeIndex) Sizes() []string {
	out := make([]string, 0, len(idx.m))
	for k := range idx.m {
		out = append(out, k)
	}
	return out
}
