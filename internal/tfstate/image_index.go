package tfstate

import "strings"

// ImageIndex maps normalised image names to the resources that use them.
type ImageIndex struct {
	index map[string][]Resource
}

// BuildImageIndex constructs an ImageIndex from the provided resources.
// Lookup is always case-insensitive (keys are lower-cased).
func BuildImageIndex(resources []Resource) *ImageIndex {
	idx := &ImageIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, _ := r.Attributes["image"].(string)
		if v == "" {
			continue
		}
		key := strings.ToLower(v)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns all resources whose image matches the given value (case-insensitive).
func (i *ImageIndex) Lookup(image string) []Resource {
	return i.index[strings.ToLower(image)]
}

// Images returns all distinct image values present in the index.
func (i *ImageIndex) Images() []string {
	out := make([]string, 0, len(i.index))
	for k := range i.index {
		out = append(out, k)
	}
	return out
}
