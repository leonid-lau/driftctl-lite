package tfstate

import "strings"

// ZoneFilterOptions configures zone filtering behaviour.
type ZoneFilterOptions struct {
	CaseSensitive bool
}

// DefaultZoneFilterOptions returns sensible defaults (case-insensitive).
func DefaultZoneFilterOptions() ZoneFilterOptions {
	return ZoneFilterOptions{CaseSensitive: false}
}

// FilterByZone returns resources whose Zone attribute matches zone.
// An empty zone string returns all resources unchanged.
func FilterByZone(resources []Resource, zone string, opts ZoneFilterOptions) []Resource {
	if zone == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["zone"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		if matchZone(s, zone, opts.CaseSensitive) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByZones returns resources matching ANY of the provided zones (OR semantics).
func FilterByZones(resources []Resource, zones []string, opts ZoneFilterOptions) []Resource {
	if len(zones) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, ok := r.Attributes["zone"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		for _, z := range zones {
			if matchZone(s, z, opts.CaseSensitive) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchZone(val, target string, caseSensitive bool) bool {
	if caseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
