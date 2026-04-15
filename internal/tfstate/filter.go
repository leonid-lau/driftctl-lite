package tfstate

// FilterByType returns only resources whose Type matches the given type string.
func FilterByType(resources []Resource, resourceType string) []Resource {
	var result []Resource
	for _, r := range resources {
		if r.Type == resourceType {
			result = append(result, r)
		}
	}
	return result
}

// IndexByID builds a map of resource ID (from attributes["id"]) to Resource
// for fast lookup during drift comparison. Resources without an "id" attribute
// are keyed by "<type>.<name>" as a fallback.
func IndexByID(resources []Resource) map[string]Resource {
	index := make(map[string]Resource, len(resources))
	for _, r := range resources {
		key := fallbackKey(r)
		if id, ok := r.Attributes["id"].(string); ok && id != "" {
			key = id
		}
		index[key] = r
	}
	return index
}

func fallbackKey(r Resource) string {
	return r.Type + "." + r.Name
}
