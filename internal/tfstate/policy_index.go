package tfstate

import "strings"

// PolicyIndex maps normalised policy values to slices of resources.
type PolicyIndex struct {
	index map[string][]Resource
}

// BuildPolicyIndex constructs a PolicyIndex from a slice of resources.
// Lookup keys are lower-cased for case-insensitive access.
func BuildPolicyIndex(resources []Resource) *PolicyIndex {
	idx := &PolicyIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["policy"]
		if !ok || v == "" {
			continue
		}
		key := strings.ToLower(v)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns all resources with the given policy (case-insensitive).
func (pi *PolicyIndex) Lookup(policy string) []Resource {
	if policy == "" {
		return nil
	}
	return pi.index[strings.ToLower(policy)]
}

// Policies returns all distinct policy values present in the index.
func (pi *PolicyIndex) Policies() []string {
	keys := make([]string, 0, len(pi.index))
	for k := range pi.index {
		keys = append(keys, k)
	}
	return keys
}
