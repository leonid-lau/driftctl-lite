package tfstate

import "strings"

// ProjectFilterOptions controls matching behaviour for project filters.
type ProjectFilterOptions struct {
	CaseSensitive bool
}

// DefaultProjectFilterOptions returns the default options (case-insensitive).
func DefaultProjectFilterOptions() ProjectFilterOptions {
	return ProjectFilterOptions{CaseSensitive: false}
}

// FilterByProject returns resources whose "project" attribute matches the
// given value. An empty project string returns all resources unchanged.
func FilterByProject(resources []Resource, project string, opts ProjectFilterOptions) []Resource {
	if project == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["project"].(string)
		if matchProject(v, project, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByProjects returns resources that match ANY of the supplied project
// values (OR semantics).
func FilterByProjects(resources []Resource, projects []string, opts ProjectFilterOptions) []Resource {
	if len(projects) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["project"].(string)
		for _, p := range projects {
			if matchProject(v, p, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchProject(val, target string, opts ProjectFilterOptions) bool {
	if opts.CaseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
