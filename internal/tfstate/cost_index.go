package tfstate

// CostIndex maps cost_center values to slices of resources.
type CostIndex struct {
	index map[string][]Resource
}

// BuildCostIndex constructs an index keyed by cost_center attribute (lowercased).
func BuildCostIndex(resources []Resource) *CostIndex {
	idx := &CostIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["cost_center"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok || s == "" {
			continue
		}
		key := toLower(s)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources for the given cost center (case-insensitive).
func (ci *CostIndex) Lookup(costCenter string) []Resource {
	if costCenter == "" {
		return nil
	}
	return ci.index[toLower(costCenter)]
}

// CostCenters returns all indexed cost center keys.
func (ci *CostIndex) CostCenters() []string {
	keys := make([]string, 0, len(ci.index))
	for k := range ci.index {
		keys = append(keys, k)
	}
	return keys
}
