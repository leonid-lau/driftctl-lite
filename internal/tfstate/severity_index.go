package tfstate

import "strings"

// SeverityIndex maps severity level -> list of resources.
type SeverityIndex struct {
	index map[string][]Resource
}

// BuildSeverityIndex constructs a SeverityIndex from a slice of resources.
func BuildSeverityIndex(resources []Resource) *SeverityIndex {
	idx := &SeverityIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["severity"]
		if !ok {
			continue
		}
		val, _ := v.(string)
		key := strings.ToLower(val)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources for the given severity (case-insensitive).
func (s *SeverityIndex) Lookup(severity string) []Resource {
	return s.index[strings.ToLower(severity)]
}

// Severities returns all indexed severity levels.
func (s *SeverityIndex) Severities() []string {
	keys := make([]string, 0, len(s.index))
	for k := range s.index {
		keys = append(keys, k)
	}
	return keys
}
