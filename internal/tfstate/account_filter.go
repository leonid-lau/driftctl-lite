package tfstate

import "strings"

// AccountFilterOptions controls matching behaviour for account filters.
type AccountFilterOptions struct {
	CaseSensitive bool
}

// DefaultAccountFilterOptions returns the default options (case-insensitive).
func DefaultAccountFilterOptions() AccountFilterOptions {
	return AccountFilterOptions{CaseSensitive: false}
}

// FilterByAccount returns resources whose account attribute matches the given
// account ID. An empty account returns all resources unchanged.
func FilterByAccount(resources []Resource, account string, opts AccountFilterOptions) []Resource {
	if account == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["account"].(string)
		if v == "" {
			v, _ = r.Attributes["account_id"].(string)
		}
		if matchAccount(v, account, opts.CaseSensitive) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByAccounts returns resources matching any of the provided account IDs
// (OR semantics).
func FilterByAccounts(resources []Resource, accounts []string, opts AccountFilterOptions) []Resource {
	if len(accounts) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["account"].(string)
		if v == "" {
			v, _ = r.Attributes["account_id"].(string)
		}
		for _, a := range accounts {
			if matchAccount(v, a, opts.CaseSensitive) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchAccount(val, target string, caseSensitive bool) bool {
	if caseSensitive {
		return val == target
	}
	return strings.EqualFold(val, target)
}
