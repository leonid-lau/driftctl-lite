package tfstate

// TagIndex provides fast lookup of resources by tag key/value pairs.
type TagIndex struct {
	// byTag maps "key=value" -> list of resource IDs
	byTag map[string][]string
	// byKey maps tag key -> list of resource IDs
	byKey map[string][]string
}

// BuildTagIndex constructs a TagIndex from a slice of Resources.
func BuildTagIndex(resources []Resource) *TagIndex {
	idx := &TagIndex{
		byTag: make(map[string][]string),
		byKey: make(map[string][]string),
	}
	for _, r := range resources {
		id := fallbackKey(r)
		for k, v := range r.Attributes {
			if len(k) > 4 && k[:5] == "tags." {
				tagKey := k[5:]
				composite := tagKey + "=" + v
				idx.byTag[composite] = append(idx.byTag[composite], id)
				idx.byKey[tagKey] = append(idx.byKey[tagKey], id)
			}
		}
	}
	return idx
}

// LookupByTag returns resource IDs matching the given tag key and value.
// If value is empty, matches any resource with the given key.
func (idx *TagIndex) LookupByTag(key, value string) []string {
	if key == "" {
		return nil
	}
	if value == "" {
		return idx.byKey[key]
	}
	return idx.byTag[key+"="+value]
}

// Keys returns all indexed tag keys.
func (idx *TagIndex) Keys() []string {
	keys := make([]string, 0, len(idx.byKey))
	for k := range idx.byKey {
		keys = append(keys, k)
	}
	return keys
}

// HasTag reports whether any indexed resource has the given tag key and value.
// If value is empty, it checks for the presence of the key with any value.
func (idx *TagIndex) HasTag(key, value string) bool {
	return len(idx.LookupByTag(key, value)) > 0
}
