package tfstate

import "github.com/owner/driftctl-lite/internal/tfstate"

// DefaultAnnotationFilterOptions returns default options for annotation filtering.
func DefaultAnnotationFilterOptions() AnnotationFilterOptions {
	return AnnotationFilterOptions{CaseSensitive: true}
}

// AnnotationFilterOptions controls annotation filter behaviour.
type AnnotationFilterOptions struct {
	CaseSensitive bool
}

// FilterByAnnotation returns resources whose annotations contain the given key
// (and optionally value). An empty key returns all resources.
func FilterByAnnotation(resources []Resource, key, value string, opts AnnotationFilterOptions) []Resource {
	if key == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		annotations, ok := r.Attributes["annotations"]
		if !ok {
			continue
		}
		annoMap, ok := annotations.(map[string]interface{})
		if !ok {
			continue
		}
		v, exists := annoMap[key]
		if !exists {
			continue
		}
		if value == "" {
			out = append(out, r)
			continue
		}
		if sv, ok := v.(string); ok && sv == value {
			out = append(out, r)
		}
	}
	return out
}

// FilterByAnnotations applies multiple key/value annotation filters with AND semantics.
func FilterByAnnotations(resources []Resource, filters map[string]string, opts AnnotationFilterOptions) []Resource {
	result := resources
	for k, v := range filters {
		result = FilterByAnnotation(result, k, v, opts)
	}
	return result
}
