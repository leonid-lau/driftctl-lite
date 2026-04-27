package tfstate

import "strings"

// SubnetIndex maps normalised subnet identifiers to slices of Resources.
type SubnetIndex struct {
	index map[string][]Resource
}

// BuildSubnetIndex constructs a SubnetIndex from the provided resources.
// Both "subnet" and "subnet_id" attributes are considered; the value is
// normalised to lower-case for case-insensitive lookups.
func BuildSubnetIndex(resources []Resource) *SubnetIndex {
	idx := &SubnetIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		key := subnetKey(r)
		if key == "" {
			continue
		}
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources associated with the given subnet (case-insensitive).
func (si *SubnetIndex) Lookup(subnet string) []Resource {
	return si.index[strings.ToLower(subnet)]
}

// Subnets returns all indexed subnet keys in no particular order.
func (si *SubnetIndex) Subnets() []string {
	keys := make([]string, 0, len(si.index))
	for k := range si.index {
		keys = append(keys, k)
	}
	return keys
}

func subnetKey(r Resource) string {
	if v, ok := r.Attributes["subnet"]; ok && v != "" {
		return strings.ToLower(v)
	}
	if v, ok := r.Attributes["subnet_id"]; ok && v != "" {
		return strings.ToLower(v)
	}
	return ""
}
