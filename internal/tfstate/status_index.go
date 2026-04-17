package tfstate

// StatusIndex maps a status string to the resources that carry it.
type StatusIndex struct {
	index map[string][]Resource
}

// BuildStatusIndex constructs a StatusIndex from a slice of resources.
func BuildStatusIndex(resources []Resource) *StatusIndex {
	idx := &StatusIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		if r.Status == "" {
			continue
		}
		idx.index[r.Status] = append(idx.index[r.Status], r)
	}
	return idx
}

// Lookup returns all resources with the given status, or nil if none.
func (si *StatusIndex) Lookup(status string) []Resource {
	if status == "" {
		return nil
	}
	return si.index[status]
}

// Statuses returns all distinct status values present in the index.
func (si *StatusIndex) Statuses() []string {
	out := make([]string, 0, len(si.index))
	for k := range si.index {
		out = append(out, k)
	}
	return out
}
