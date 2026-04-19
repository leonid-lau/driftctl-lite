package tfstate

import "strings"

// DefaultModuleFilterOptions returns case-insensitive matching by default.
type ModuleFilterOptions struct {
	CaseSensitive bool
}

func DefaultModuleFilterOptions() ModuleFilterOptions {
	return ModuleFilterOptions{CaseSensitive: false}
}

// FilterByModule returns resources whose module attribute matches the given value.
// An empty module string returns all resources unchanged.
func FilterByModule(resources []Resource, module string, opts ModuleFilterOptions) []Resource {
	if module == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["module"].(string)
		if matchModule(v, module, opts) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByModules returns resources matching ANY of the given module values (OR semantics).
func FilterByModules(resources []Resource, modules []string, opts ModuleFilterOptions) []Resource {
	if len(modules) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["module"].(string)
		for _, m := range modules {
			if matchModule(v, m, opts) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchModule(val, target string, opts ModuleFilterOptions) bool {
	if opts.CaseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
