package tfstate

import "strings"

// Severity levels for drift classification.
const (
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"
)

// DefaultSeverityFilterOptions returns case-insensitive matching.
type SeverityFilterOptions struct {
	CaseInsensitive bool
}

func DefaultSeverityFilterOptions() SeverityFilterOptions {
	return SeverityFilterOptions{CaseInsensitive: true}
}

// FilterBySeverity returns resources whose Severity attribute matches the given level.
func FilterBySeverity(resources []Resource, severity string, opts SeverityFilterOptions) []Resource {
	if severity == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["severity"]
		if !ok {
			continue
		}
		val, _ := v.(string)
		if opts.CaseInsensitive {
			if strings.EqualFold(val, severity) {
				out = append(out, r)
			}
		} else {
			if val == severity {
				out = append(out, r)
			}
		}
	}
	return out
}

// FilterBySeverities returns resources matching any of the given severity levels (OR semantics).
func FilterBySeverities(resources []Resource, severities []string, opts SeverityFilterOptions) []Resource {
	if len(severities) == 0 {
		return resources
	}
	set := make(map[string]struct{}, len(severities))
	for _, s := range severities {
		key := s
		if opts.CaseInsensitive {
			key = strings.ToLower(s)
		}
		set[key] = struct{}{}
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["severity"]
		if !ok {
			continue
		}
		val, _ := v.(string)
		key := val
		if opts.CaseInsensitive {
			key = strings.ToLower(val)
		}
		if _, found := set[key]; found {
			out = append(out, r)
		}
	}
	return out
}
