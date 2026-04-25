package tfstate

import "strings"

// SKUIndex maps normalised SKU values to slices of resources.
type SKUIndex struct {
	index map[string][]Resource
}

// BuildSKUIndex constructs a case-insensitive index over the "sku" attribute.
func BuildSKUIndex(resources []Resource) *SKUIndex {
	idx := &SKUIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		val, ok := r.Attributes["sku"]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok || s == "" {
			continue
		}
		key := strings.ToLower(s)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns all resources with the given SKU (case-insensitive).
func (i *SKUIndex) Lookup(sku string) []Resource {
	if sku == "" {
		return nil
	}
	return i.index[strings.ToLower(sku)]
}

// SKUs returns all distinct SKU values present in the index.
func (i *SKUIndex) SKUs() []string {
	out := make([]string, 0, len(i.index))
	for k := range i.index {
		out = append(out, k)
	}
	return out
}
