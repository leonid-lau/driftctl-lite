package tfstate

import "strings"

// TierIndex maps normalised tier value -> resources.
type TierIndex struct {
	index map[string][]Resource
}

// BuildTierIndex constructs a TierIndex from a slice of resources.
func BuildTierIndex(resources []Resource) *TierIndex {
	idx := &TierIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["tier"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		key := strings.ToLower(s)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources matching the given tier (case-insensitive).
func (ti *TierIndex) Lookup(tier string) []Resource {
	if tier == "" {
		return nil
	}
	return ti.index[strings.ToLower(tier)]
}

// Tiers returns all distinct tier values present in the index.
func (ti *TierIndex) Tiers() []string {
	keys := make([]string, 0, len(ti.index))
	for k := range ti.index {
		keys = append(keys, k)
	}
	return keys
}
