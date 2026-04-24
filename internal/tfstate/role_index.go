package tfstate

import "strings"

// RoleIndex maps normalised role values to the resources that carry them.
type RoleIndex struct {
	index map[string][]Resource
}

// BuildRoleIndex constructs a RoleIndex from the provided resources.
// Keys are stored in lower-case to enable case-insensitive look-ups.
func BuildRoleIndex(resources []Resource) *RoleIndex {
	idx := &RoleIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		val, ok := r.Attributes["role"]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok || s == "" {
			continue
		}
		key := strings.ToLower(s)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns all resources with the given role (case-insensitive).
// Returns nil when no resources match.
func (ri *RoleIndex) Lookup(role string) []Resource {
	if role == "" {
		return nil
	}
	return ri.index[strings.ToLower(role)]
}

// Roles returns the distinct (lower-cased) role values present in the index.
func (ri *RoleIndex) Roles() []string {
	keys := make([]string, 0, len(ri.index))
	for k := range ri.index {
		keys = append(keys, k)
	}
	return keys
}
