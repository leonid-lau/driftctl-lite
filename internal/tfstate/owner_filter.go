package tfstate

import "github.com/owner/driftctl-lite/internal/tfstate"

// DefaultOwnerFilterOptions returns sensible defaults.
func DefaultOwnerFilterOptions() OwnerFilterOptions {
	return OwnerFilterOptions{CaseSensitive: false}
}

// OwnerFilterOptions controls matching behaviour.
type OwnerFilterOptions struct {
	CaseSensitive bool
}

// FilterByOwner returns resources whose Owner metadata field matches owner.
// An empty owner string returns all resources unchanged.
func FilterByOwner(resources []Resource, owner string, opts OwnerFilterOptions) []Resource {
	if owner == "" {
		return resources
	}
	want := owner
	if !opts.CaseSensitive {
		want = toLower(owner)
	}
	var out []Resource
	for _, r := range resources {
		val := r.Metadata["owner"]
		if !opts.CaseSensitive {
			val = toLower(val)
		}
		if val == want {
			out = append(out, r)
		}
	}
	return out
}

// FilterByOwners returns resources matching ANY of the given owners (OR semantics).
func FilterByOwners(resources []Resource, owners []string, opts OwnerFilterOptions) []Resource {
	if len(owners) == 0 {
		return resources
	}
	wantSet := make(map[string]struct{}, len(owners))
	for _, o := range owners {
		key := o
		if !opts.CaseSensitive {
			key = toLower(o)
		}
		wantSet[key] = struct{}{}
	}
	var out []Resource
	for _, r := range resources {
		val := r.Metadata["owner"]
		if !opts.CaseSensitive {
			val = toLower(val)
		}
		if _, ok := wantSet[val]; ok {
			out = append(out, r)
		}
	}
	return out
}
