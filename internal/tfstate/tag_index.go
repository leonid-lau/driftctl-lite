package tfstate

// TagIndex maps tag key -> tag value -> list of resources.
type TagIndex struct {
	idx map[string]map[string][]Resource
}

// BuildTagIndex constructs a TagIndex from the given resources.
func BuildTagIndex(resources []Resource) *TagIndex {
	idx := &TagIndex{idx: make(map[string]map[string][]Resource)}
	for _, r := range resources {
		for attrKey, attrVal := range r.Attributes {
			var tagKey string
			if len(attrKey) > 5 && attrKey[:5] == "tags." {
				tagKey = attrKey[5:]
			} else if len(attrKey) > 4 && attrKey[:4] == "tag." {
				tagKey = attrKey[4:]
			} else {
				continue
			}
			if idx.idx[tagKey] == nil {
				idx.idx[tagKey] = make(map[string][]Resource)
			}
			idx.idx[tagKey][attrVal] = append(idx.idx[tagKey][attrVal], r)
		}
	}
	return idx
}

// Lookup returns resources matching the key (and optionally value).
// If value is empty, all resources with that tag key are returned.
func (t *TagIndex) Lookup(key, value string) []Resource {
	if key == "" {
		return nil
	}
	valMap, ok := t.idx[key]
	if !ok {
		return nil
	}
	if value == "" {
		var all []Resource
		for _, rs := range valMap {
			all = append(all, rs...)
		}
		return all
	}
	return valMap[value]
}

// Keys returns all indexed tag keys.
func (t *TagIndex) Keys() []string {
	keys := make([]string, 0, len(t.idx))
	for k := range t.idx {
		keys = append(keys, k)
	}
	return keys
}
