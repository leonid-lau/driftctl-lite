package tfstate

// LabelIndex maps label key -> value -> list of resources.
type LabelIndex struct {
	index map[string]map[string][]Resource
}

// BuildLabelIndex constructs a LabelIndex from a slice of resources.
func BuildLabelIndex(resources []Resource) *LabelIndex {
	idx := &LabelIndex{index: make(map[string]map[string][]Resource)}
	for _, r := range resources {
		for k, v := range r.Labels {
			if idx.index[k] == nil {
				idx.index[k] = make(map[string][]Resource)
			}
			idx.index[k][v] = append(idx.index[k][v], r)
		}
	}
	return idx
}

// Lookup returns resources matching the given label key and value.
// If value is empty, all resources with that key are returned.
func (li *LabelIndex) Lookup(key, value string) []Resource {
	if key == "" {
		return nil
	}
	values, ok := li.index[key]
	if !ok {
		return nil
	}
	if value == "" {
		var out []Resource
		for _, rs := range values {
			out = append(out, rs...)
		}
		return out
	}
	return values[value]
}

// Keys returns all label keys present in the index.
func (li *LabelIndex) Keys() []string {
	keys := make([]string, 0, len(li.index))
	for k := range li.index {
		keys = append(keys, k)
	}
	return keys
}
