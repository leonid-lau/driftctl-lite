package tfstate

import "strings"

// DefaultScheduleFilterOptions returns options with case-insensitive matching enabled.
func DefaultScheduleFilterOptions() FilterOptions {
	return FilterOptions{CaseInsensitive: true}
}

// FilterOptions holds configuration for filter operations.
type FilterOptions struct {
	CaseInsensitive bool
}

// FilterBySchedule returns resources whose "schedule" attribute matches the given value.
// An empty schedule returns all resources unchanged.
func FilterBySchedule(resources []Resource, schedule string, opts FilterOptions) []Resource {
	if schedule == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		if matchSchedule(r, schedule, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterBySchedules returns resources that match ANY of the provided schedules (OR semantics).
func FilterBySchedules(resources []Resource, schedules []string, opts FilterOptions) []Resource {
	if len(schedules) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		for _, s := range schedules {
			if matchSchedule(r, s, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchSchedule(r Resource, schedule string, opts FilterOptions) bool {
	v, ok := r.Attributes["schedule"]
	if !ok {
		return false
	}
	val, _ := v.(string)
	if opts.CaseInsensitive {
		return strings.EqualFold(val, schedule)
	}
	return val == schedule
}
