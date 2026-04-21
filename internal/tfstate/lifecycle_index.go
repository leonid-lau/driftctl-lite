package tfstate

import "strings"

// LifecycleIndex maps normalised lifecycle value -> list of resources.
type LifecycleIndex struct {
	m map[string][]Resource
}

// BuildLifecycleIndex builds an index from a slice of resources.
// Keys are stored in lower-case for case-insensitive lookups.
func BuildLifecycleIndex(resources []Resource) *LifecycleIndex {
	m := make(map[string][]Resource)
	for _, r := range resources {
		v, ok := r.Attributes["lifecycle"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok || s == "" {
			continue
		}
		key := strings.ToLower(s)
		m[key] = append(m[key], r)
	}
	return &LifecycleIndex{m: m}
}

// Lookup returns resources with the given lifecycle (case-insensitive).
func (idx *LifecycleIndex) Lookup(lifecycle string) []Resource {
	if lifecycle == "" {
		return nil
	}
	return idx.m[strings.ToLower(lifecycle)]
}

// Lifecycles returns all distinct lifecycle values present in the index.
func (idx *LifecycleIndex) Lifecycles() []string {
	keys := make([]string, 0, len(idx.m))
	for k := range idx.m {
		keys = append(keys, k)
	}
	return keys
}
