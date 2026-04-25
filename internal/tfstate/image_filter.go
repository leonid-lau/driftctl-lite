package tfstate

import "strings"

// ImageFilterOptions controls matching behaviour for image filters.
type ImageFilterOptions struct {
	CaseSensitive bool
}

// DefaultImageFilterOptions returns the default options (case-insensitive).
func DefaultImageFilterOptions() ImageFilterOptions {
	return ImageFilterOptions{CaseSensitive: false}
}

// FilterByImage returns resources whose "image" attribute matches the given value.
// An empty image string returns all resources unchanged.
func FilterByImage(resources []Resource, image string, opts ImageFilterOptions) []Resource {
	if image == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchImage(r, image, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByImages returns resources that match ANY of the provided image values (OR semantics).
func FilterByImages(resources []Resource, images []string, opts ImageFilterOptions) []Resource {
	if len(images) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, img := range images {
			if matchImage(r, img, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchImage(r Resource, image string, opts ImageFilterOptions) bool {
	v, _ := r.Attributes["image"].(string)
	if opts.CaseSensitive {
		return v == image
	}
	return strings.EqualFold(v, image)
}
