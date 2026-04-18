package tfstate

import "strings"

// Priority levels in ascending order of importance.
const (
	PriorityLow      = "low"
	PriorityMedium   = "medium"
	PriorityHigh     = "high"
	PriorityCritical = "critical"
)

var priorityRank = map[string]int{
	PriorityLow:      0,
	PriorityMedium:   1,
	PriorityHigh:     2,
	PriorityCritical: 3,
}

// DefaultPriorityFilterOptions returns case-insensitive matching options.
type PriorityFilterOptions struct {
	CaseSensitive bool
	MinPriority   string // if set, include resources at or above this level
}

func DefaultPriorityFilterOptions() PriorityFilterOptions {
	return PriorityFilterOptions{CaseSensitive: false}
}

// FilterByPriority returns resources whose Priority attribute matches priority.
// An empty priority string returns all resources.
func FilterByPriority(resources []Resource, priority string, opts PriorityFilterOptions) []Resource {
	if priority == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["priority"].(string)
		if match(v, priority, opts.CaseSensitive) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByMinPriority returns resources whose priority rank is >= the given level.
func FilterByMinPriority(resources []Resource, minPriority string) []Resource {
	if minPriority == "" {
		return resources
	}
	minRank, ok := priorityRank[strings.ToLower(minPriority)]
	if !ok {
		return nil
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["priority"].(string)
		rank, known := priorityRank[strings.ToLower(v)]
		if known && rank >= minRank {
			out = append(out, r)
		}
	}
	return out
}

func match(value, target string, caseSensitive bool) bool {
	if caseSensitive {
		return value == target
	}
	return strings.EqualFold(value, target)
}
