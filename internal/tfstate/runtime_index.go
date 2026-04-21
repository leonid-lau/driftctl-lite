package tfstate

// BuildRuntimeIndex builds an index from runtime value to resources.
// Keys are always stored lower-case for case-insensitive lookup.
func BuildRuntimeIndex(resources []Resource) map[string][]Resource {
	idx := make(map[string][]Resource)
	for _, r := range resources {
		val, ok := r.Attributes["runtime"]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		key := toLower(s)
		idx[key] = append(idx[key], r)
	}
	return idx
}

// LookupByRuntime returns resources matching the given runtime (case-insensitive).
func LookupByRuntime(idx map[string][]Resource, runtime string) []Resource {
	return idx[toLower(runtime)]
}

// Runtimes returns all distinct runtime values present in the index.
func Runtimes(idx map[string][]Resource) []string {
	keys := make([]string, 0, len(idx))
	for k := range idx {
		keys = append(keys, k)
	}
	return keys
}
