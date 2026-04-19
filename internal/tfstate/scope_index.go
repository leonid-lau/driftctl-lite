package tfstate

import "strings"

// ScopeIndex maps normalised scope values to the resources that carry them.
type ScopeIndex struct {
	idx map[string][]Resource
}

// BuildScopeIndex constructs a ScopeIndex from a slice of resources.
// Lookup keys are lower-cased for case-insensitive matching.
func BuildScopeIndex(resources []Resource) *ScopeIndex {
	idx := make(map[string][]Resource)
	for _, r := range resources {
		if v, ok := r.Attributes["scope"]; ok && v != "" {
			key := strings.ToLower(v)
			idx[key] = append(idx[key], r)
		}
	}
	return &ScopeIndex{idx: idx}
}

// Lookup returns resources with the given scope (case-insensitive).
func (si *ScopeIndex) Lookup(scope string) []Resource {
	if scope == "" {
		return nil
	}
	return si.idx[strings.ToLower(scope)]
}

// Scopes returns all indexed scope values (normalised to lower-case).
func (si *ScopeIndex) Scopes() []string {
	out := make([]string, 0, len(si.idx))
	for k := range si.idx {
		out = append(out, k)
	}
	return out
}
