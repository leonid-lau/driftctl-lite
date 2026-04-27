package tfstate

import "strings"

// ScheduleIndex maps normalised schedule values to the resources that carry them.
type ScheduleIndex struct {
	index map[string][]Resource
}

// BuildScheduleIndex constructs a ScheduleIndex from the provided resources.
// Keys are stored in lower-case to allow case-insensitive lookup.
func BuildScheduleIndex(resources []Resource) *ScheduleIndex {
	idx := &ScheduleIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, ok := r.Attributes["schedule"]
		if !ok {
			continue
		}
		val, _ := v.(string)
		if val == "" {
			continue
		}
		key := strings.ToLower(val)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources matching the given schedule (case-insensitive).
// Returns nil when no resources are found.
func (si *ScheduleIndex) Lookup(schedule string) []Resource {
	if schedule == "" {
		return nil
	}
	return si.index[strings.ToLower(schedule)]
}

// Schedules returns all distinct schedule values present in the index.
func (si *ScheduleIndex) Schedules() []string {
	keys := make([]string, 0, len(si.index))
	for k := range si.index {
		keys = append(keys, k)
	}
	return keys
}
