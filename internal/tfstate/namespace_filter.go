package tfstate

// NamespaceFilterOptions configures namespace filtering behaviour.
type NamespaceFilterOptions struct {
	// CaseSensitive controls whether namespace matching is case-sensitive.
	CaseSensitive bool
}

// DefaultNamespaceFilterOptions returns sensible defaults.
func DefaultNamespaceFilterOptions() NamespaceFilterOptions {
	return NamespaceFilterOptions{CaseSensitive: true}
}

// FilterByNamespace returns resources whose Namespace field equals ns.
// An empty ns returns all resources unchanged.
func FilterByNamespace(resources []Resource, ns string, opts NamespaceFilterOptions) []Resource {
	if ns == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		rns := r.Namespace
		cmp := ns
		if !opts.CaseSensitive {
			rns = toLower(rns)
			cmp = toLower(cmp)
		}
		if rns == cmp {
			out = append(out, r)
		}
	}
	return out
}

// FilterByNamespaces returns resources matching ANY of the provided namespaces.
func FilterByNamespaces(resources []Resource, namespaces []string, opts NamespaceFilterOptions) []Resource {
	if len(namespaces) == 0 {
		return resources
	}
	set := make(map[string]struct{}, len(namespaces))
	for _, ns := range namespaces {
		k := ns
		if !opts.CaseSensitive {
			k = toLower(k)
		}
		set[k] = struct{}{}
	}
	var out []Resource
	for _, r := range resources {
		k := r.Namespace
		if !opts.CaseSensitive {
			k = toLower(k)
		}
		if _, ok := set[k]; ok {
			out = append(out, r)
		}
	}
	return out
}

func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
