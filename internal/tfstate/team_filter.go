package tfstate

// DefaultTeamFilterOptions returns a TeamFilterOptions with CaseSensitive=false.
func DefaultTeamFilterOptions() TeamFilterOptions {
	return TeamFilterOptions{CaseSensitive: false}
}

// TeamFilterOptions controls matching behaviour for team filters.
type TeamFilterOptions struct {
	CaseSensitive bool
}

// FilterByTeam returns resources whose "team" attribute matches team.
// An empty team returns all resources unchanged.
func FilterByTeam(resources []Resource, team string, opts TeamFilterOptions) []Resource {
	if team == "" {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["team"].(string)
		if matchTeam(v, team, opts.CaseSensitive) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByTeams returns resources matching ANY of the given teams (OR semantics).
func FilterByTeams(resources []Resource, teams []string, opts TeamFilterOptions) []Resource {
	if len(teams) == 0 {
		return resources
	}
	var out []Resource
	for _, r := range resources {
		v, _ := r.Attributes["team"].(string)
		for _, t := range teams {
			if matchTeam(v, t, opts.CaseSensitive) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func matchTeam(val, target string, caseSensitive bool) bool {
	if caseSensitive {
		return val == target
	}
	return toLower(val) == toLower(target)
}
