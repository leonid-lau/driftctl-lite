package tfstate

import "strings"

// FormatIndex maps normalised format values to slices of resources.
type FormatIndex struct {
	index map[string][]Resource
}

// BuildFormatIndex constructs a FormatIndex from the provided resources.
// Keys are stored in lower-case for case-insensitive lookup.
func BuildFormatIndex(resources []Resource) *FormatIndex {
	idx := &FormatIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["format"]
		if !ok || v == "" {
			continue
		}
		key := strings.ToLower(v)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources whose format matches the given value (case-insensitive).
func (fi *FormatIndex) Lookup(format string) []Resource {
	if format == "" {
		return nil
	}
	return fi.index[strings.ToLower(format)]
}

// Formats returns all distinct format values present in the index.
func (fi *FormatIndex) Formats() []string {
	keys := make([]string, 0, len(fi.index))
	for k := range fi.index {
		keys = append(keys, k)
	}
	return keys
}
