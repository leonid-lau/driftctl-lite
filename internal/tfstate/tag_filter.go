package tfstate

import "github.com/driftctl-lite/internal/tfstate"

// TagFilter holds criteria for filtering resources by their tags.
type TagFilter struct {
	Key   string
	Value string // empty means match any value for the key
}

// FilterByTag returns resources whose attributes contain a matching tag entry.
// Tags are expected under attributes as "tag.<key>" = "<value>".
func FilterByTag(resources []Resource, f TagFilter) []Resource {
	if f.Key == "" {
		return resources
	}
	tagKey := "tag." + f.Key
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes[tagKey]
		if !ok {
			continue
		}
		if f.Value == "" || v == f.Value {
			out = append(out, r)
		}
	}
	return out
}

// FilterByTags applies multiple TagFilters (AND semantics).
func FilterByTags(resources []Resource, filters []TagFilter) []Resource {
	result := resources
	for _, f := range filters {
		result = FilterByTag(result, f)
	}
	return result
}
