package tfstate

// TeamIndex maps team name (lower-case) to resources.
type TeamIndex struct {
	index map[string][]Resource
}

// BuildTeamIndex constructs a TeamIndex from a slice of resources.
func BuildTeamIndex(resources []Resource) *TeamIndex {
	idx := &TeamIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		v, _ := r.Attributes["team"].(string)
		if v == "" {
			continue
		}
		key := toLower(v)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns resources belonging to team (case-insensitive).
func (i *TeamIndex) Lookup(team string) []Resource {
	if team == "" {
		return nil
	}
	return i.index[toLower(team)]
}

// Teams returns all known team names stored in the index.
func (i *TeamIndex) Teams() []string {
	keys := make([]string, 0, len(i.index))
	for k := range i.index {
		keys = append(keys, k)
	}
	return keys
}
