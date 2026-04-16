package tfstate

// AnnotationIndex provides fast lookup of resources by annotation key/value.
type AnnotationIndex struct {
	// index[key][value] -> list of resource indices
	index map[string]map[string][]int
	resources []Resource
}

// BuildAnnotationIndex constructs an AnnotationIndex from a slice of resources.
func BuildAnnotationIndex(resources []Resource) *AnnotationIndex {
	idx := &AnnotationIndex{
		index:     make(map[string]map[string][]int),
		resources: resources,
	}
	for i, r := range resources {
		for k, v := range r.Annotations {
			if _, ok := idx.index[k]; !ok {
				idx.index[k] = make(map[string][]int)
			}
			idx.index[k][v] = append(idx.index[k][v], i)
		}
	}
	return idx
}

// Lookup returns resources matching the given annotation key and optional value.
// If value is empty, all resources with the key are returned.
func (a *AnnotationIndex) Lookup(key, value string) []Resource {
	if key == "" {
		return nil
	}
	valMap, ok := a.index[key]
	if !ok {
		return nil
	}
	var indices []int
	if value == "" {
		for _, idxList := range valMap {
			indices = append(indices, idxList...)
		}
	} else {
		indices = valMap[value]
	}
	result := make([]Resource, 0, len(indices))
	for _, i := range indices {
		result = append(result, a.resources[i])
	}
	return result
}

// Keys returns all annotation keys present in the index.
func (a *AnnotationIndex) Keys() []string {
	keys := make([]string, 0, len(a.index))
	for k := range a.index {
		keys = append(keys, k)
	}
	return keys
}
