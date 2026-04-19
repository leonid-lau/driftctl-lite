package tfstate

import "strings"

// WorkspaceIndex maps lowercase workspace name -> resources.
type WorkspaceIndex struct {
	m map[string][]Resource
}

// BuildWorkspaceIndex builds an index keyed by the "workspace" attribute (lowercased).
func BuildWorkspaceIndex(resources []Resource) *WorkspaceIndex {
	m := make(map[string][]Resource)
	for _, r := range resources {
		v, ok := r.Attributes["workspace"]
		if !ok {
			continue
		}
		s, _ := v.(string)
		key := strings.ToLower(s)
		m[key] = append(m[key], r)
	}
	return &WorkspaceIndex{m: m}
}

// Lookup returns resources for the given workspace (case-insensitive).
func (idx *WorkspaceIndex) Lookup(ws string) []Resource {
	return idx.m[strings.ToLower(ws)]
}

// Workspaces returns all known workspace names (lowercased).
func (idx *WorkspaceIndex) Workspaces() []string {
	keys := make([]string, 0, len(idx.m))
	for k := range idx.m {
		keys = append(keys, k)
	}
	return keys
}
