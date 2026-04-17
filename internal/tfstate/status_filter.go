package tfstate

import "strings"

// ResourceStatus represents the lifecycle status of a resource.
type ResourceStatus string

const (
	StatusCreated  ResourceStatus = "created"
	StatusDestroyed ResourceStatus = "destroyed"
	StatusTainted  ResourceStatus = "tainted"
)

// DefaultStatusFilterOptions returns options with case-insensitive matching enabled.
type StatusFilterOptions struct {
	CaseInsensitive bool
}

func DefaultStatusFilterOptions() StatusFilterOptions {
	return StatusFilterOptions{CaseInsensitive: true}
}

// FilterByStatus returns resources whose Status field matches the given status.
func FilterByStatus(resources []Resource, status ResourceStatus, opts StatusFilterOptions) []Resource {
	if status == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		rs := r.Status
		st := string(status)
		if opts.CaseInsensitive {
			rs = strings.ToLower(rs)
			st = strings.ToLower(st)
		}
		if rs == st {
			out = append(out, r)
		}
	}
	return out
}

// FilterByStatuses returns resources matching any of the provided statuses (OR semantics).
func FilterByStatuses(resources []Resource, statuses []ResourceStatus, opts StatusFilterOptions) []Resource {
	if len(statuses) == 0 {
		return resources
	}
	set := make(map[string]struct{}, len(statuses))
	for _, s := range statuses {
		key := string(s)
		if opts.CaseInsensitive {
			key = strings.ToLower(key)
		}
		set[key] = struct{}{}
	}
	var out []Resource
	for _, r := range resources {
		rs := r.Status
		if opts.CaseInsensitive {
			rs = strings.ToLower(rs)
		}
		if _, ok := set[rs]; ok {
			out = append(out, r)
		}
	}
	return out
}
