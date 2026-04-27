package tfstate

// FilterByTag returns resources matching the given tag key (and optionally value).
// If key is empty, all resources are returned.
// If value is empty, any resource with the key present is returned.
func FilterByTag(resources []Resource, key, value string) []Resource {
	if key == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["tags." + key]
		if !ok {
			v, ok = r.Attributes["tag." + key]
		}
		if !ok {
			continue
		}
		if value == "" || v == value {
			out = append(out, r)
		}
	}
	return out
}

// FilterByTags applies AND semantics across all provided key/value pairs.
// A resource must match ALL pairs to be included.
func FilterByTags(resources []Resource, tags map[string]string) []Resource {
	out := resources
	for k, v := range tags {
		out = FilterByTag(out, k, v)
	}
	return out
}
