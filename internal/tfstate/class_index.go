package tfstate

import "strings"

// ClassIndex maps normalised class values to the resources that carry them.
type ClassIndex struct {
	index map[string][]Resource
}

// BuildClassIndex constructs a ClassIndex from the provided resources.
// Keys are stored in lower-case for case-insensitive look-up.
func BuildClassIndex(resources []Resource) *ClassIndex {
	idx := &ClassIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["class"]
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

// Lookup returns resources with the given class (case-insensitive).
// Returns nil when no match is found.
func (i *ClassIndex) Lookup(class string) []Resource {
	if class == "" {
		return nil
	}
	return i.index[strings.ToLower(class)]
}

// Classes returns all distinct class values present in the index.
func (i *ClassIndex) Classes() []string {
	out := make([]string, 0, len(i.index))
	for k := range i.index {
		out = append(out, k)
	}
	return out
}
