package tfstate

// ProjectIndex maps project names to matching resources.
type ProjectIndex struct {
	index map[string][]Resource
}

// BuildProjectIndex builds an index of resources keyed by their project attribute.
// Lookup is case-insensitive.
func BuildProjectIndex(resources []Resource) *ProjectIndex {
	idx := make(map[string][]Resource)
	for _, r := range resources {
		project := ""
		if v, ok := r.Attributes["project"]; ok {
			project = toLower(v)
		}
		if project == "" {
			continue
		}
		idx[project] = append(idx[project], r)
	}
	return &ProjectIndex{index: idx}
}

// Lookup returns resources matching the given project name (case-insensitive).
func (i *ProjectIndex) Lookup(project string) []Resource {
	if project == "" {
		return nil
	}
	return i.index[toLower(project)]
}

// Projects returns all indexed project names in their normalised (lowercase) form.
func (i *ProjectIndex) Projects() []string {
	keys := make([]string, 0, len(i.index))
	for k := range i.index {
		keys = append(keys, k)
	}
	return keys
}
