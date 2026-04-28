package tfstate

import "strings"

// ProtocolIndex maps protocol values to matching resources.
type ProtocolIndex struct {
	index map[string][]Resource
}

// BuildProtocolIndex constructs a ProtocolIndex from a slice of resources.
// Lookup keys are stored in lower-case for case-insensitive queries.
func BuildProtocolIndex(resources []Resource) *ProtocolIndex {
	idx := &ProtocolIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		protocol := strings.TrimSpace(r.Attributes["protocol"])
		if protocol == "" {
			continue
		}
		key := strings.ToLower(protocol)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns all resources whose protocol matches the given value
// (case-insensitive). Returns nil when no match is found.
func (idx *ProtocolIndex) Lookup(protocol string) []Resource {
	if protocol == "" {
		return nil
	}
	return idx.index[strings.ToLower(protocol)]
}

// Protocols returns the sorted list of distinct protocol keys present in the
// index.
func (idx *ProtocolIndex) Protocols() []string {
	keys := make([]string, 0, len(idx.index))
	for k := range idx.index {
		keys = append(keys, k)
	}
	return keys
}
