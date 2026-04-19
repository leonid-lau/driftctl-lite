package tfstate

import "strings"

// DefaultWorkspaceFilterOptions returns case-insensitive matching by default.
type WorkspaceFilterOptions struct {
	CaseSensitive bool
}

func DefaultWorkspaceFilterOptions() WorkspaceFilterOptions {
	return WorkspaceFilterOptions{CaseSensitive: false}
}

// FilterByWorkspace returns resources whose "workspace" attribute matches ws.
// If ws is empty, all resources are returned.
func FilterByWorkspace(resources []Resource, ws string, opts WorkspaceFilterOptions) []Resource {
	if ws == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["workspace"]
		if !ok {
			continue
		}
		s, _ := v.(string)
		if matchWorkspace(s, ws, opts.CaseSensitive) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByWorkspaces returns resources matching any of the given workspaces (OR semantics).
func FilterByWorkspaces(resources []Resource, workspaces []string, opts WorkspaceFilterOptions) []Resource {
	if len(workspaces) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["workspace"]
		if !ok {
			continue
		}
		s, _ := v.(string)
		for _, ws := range workspaces {
			if matchWorkspace(s, ws, opts.CaseSensitive) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchWorkspace(val, target string, caseSensitive bool) bool {
	if caseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
