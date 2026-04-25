package tfstate

import "strings"

// GroupIndex maps normalised group names to the resources that belong to them.
type GroupIndex struct {
	index map[string][]Resource
}

// BuildGroupIndex constructs a GroupIndex from the provided resources.
// Lookup keys are stored in lower-case for case-insensitive access.
func BuildGroupIndex(resources []Resource) *GroupIndex {
	idx := &GroupIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["group"]
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

// Lookup returns all resources belonging to the given group (case-insensitive).
// Returns nil when the group is not found.
func (gi *GroupIndex) Lookup(group string) []Resource {
	if group == "" {
		return nil
	}
	return gi.index[strings.ToLower(group)]
}

// Groups returns all distinct group names present in the index.
func (gi *GroupIndex) Groups() []string {
	keys := make([]string, 0, len(gi.index))
	for k := range gi.index {
		keys = append(keys, k)
	}
	return keys
}
