package tfstate

import "strings"

// ServiceIndex maps lowercase service name -> list of resources.
type ServiceIndex struct {
	index map[string][]Resource
}

// BuildServiceIndex builds an index of resources keyed by their "service" attribute.
func BuildServiceIndex(resources []Resource) *ServiceIndex {
	idx := &ServiceIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["service"]
		if !ok {
			continue
		}
		s, _ := v.(string)
		key := strings.ToLower(s)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns all resources for the given service (case-insensitive).
func (si *ServiceIndex) Lookup(service string) []Resource {
	return si.index[strings.ToLower(service)]
}

// Services returns all indexed service names (lowercase).
func (si *ServiceIndex) Services() []string {
	keys := make([]string, 0, len(si.index))
	for k := range si.index {
		keys = append(keys, k)
	}
	return keys
}
