package tfstate

import "strings"

// NetworkIndex maps normalised network names to the resources that belong to them.
type NetworkIndex struct {
	index map[string][]Resource
}

// BuildNetworkIndex constructs a NetworkIndex from a slice of resources.
// Lookup keys are lower-cased for case-insensitive access.
func BuildNetworkIndex(resources []Resource) *NetworkIndex {
	idx := &NetworkIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		key := networkKey(r)
		if key == "" {
			continue
		}
		k := strings.ToLower(key)
		idx.index[k] = append(idx.index[k], r)
	}
	return idx
}

// Lookup returns all resources belonging to the given network (case-insensitive).
func (i *NetworkIndex) Lookup(network string) []Resource {
	return i.index[strings.ToLower(network)]
}

// Networks returns all distinct network names stored in the index.
func (i *NetworkIndex) Networks() []string {
	keys := make([]string, 0, len(i.index))
	for k := range i.index {
		keys = append(keys, k)
	}
	return keys
}

func networkKey(r Resource) string {
	if v, ok := r.Attributes["network"]; ok {
		if s, _ := v.(string); s != "" {
			return s
		}
	}
	if v, ok := r.Attributes["vpc"]; ok {
		if s, _ := v.(string); s != "" {
			return s
		}
	}
	return ""
}
