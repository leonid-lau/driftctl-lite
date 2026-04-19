package tfstate

// AccountIndex maps account identifiers to resources.
type AccountIndex struct {
	byAccount map[string][]Resource
}

// BuildAccountIndex constructs an AccountIndex from a slice of resources.
// It normalises keys to lower-case so look-ups are case-insensitive.
func BuildAccountIndex(resources []Resource) *AccountIndex {
	idx := &AccountIndex{byAccount: make(map[string][]Resource)}
	for _, r := range resources {
		key := accountKey(r)
		if key == "" {
			continue
		}
		norm := toLower(key)
		idx.byAccount[norm] = append(idx.byAccount[norm], r)
	}
	return idx
}

// Lookup returns all resources matching the given account (case-insensitive).
func (i *AccountIndex) Lookup(account string) []Resource {
	if account == "" {
		return nil
	}
	return i.byAccount[toLower(account)]
}

// Accounts returns all distinct account keys present in the index.
func (i *AccountIndex) Accounts() []string {
	keys := make([]string, 0, len(i.byAccount))
	for k := range i.byAccount {
		keys = append(keys, k)
	}
	return keys
}

func accountKey(r Resource) string {
	if v, ok := r.Attributes["account"]; ok && v != "" {
		return v
	}
	if v, ok := r.Attributes["account_id"]; ok && v != "" {
		return v
	}
	return ""
}
