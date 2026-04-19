package tfstate

import (
	"testing"
)

func teamResource(id, team string) Resource {
	return Resource{
		Type: "aws_instance",
		ID:   id,
		Attributes: map[string]interface{}{"team": team},
	}
}

func TestFilterByTeam_Match(t *testing.T) {
	resources := []Resource{
		teamResource("r1", "platform"),
		teamResource("r2", "security"),
		teamResource("r3", "platform"),
	}
	got := FilterByTeam(resources, "platform", DefaultTeamFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByTeam_EmptyTeam_ReturnsAll(t *testing.T) {
	resources := []Resource{teamResource("r1", "platform"), teamResource("r2", "security")}
	got := FilterByTeam(resources, "", DefaultTeamFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByTeam_CaseInsensitive(t *testing.T) {
	resources := []Resource{teamResource("r1", "Platform")}
	got := FilterByTeam(resources, "platform", DefaultTeamFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByTeam_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{teamResource("r1", "Platform")}
	got := FilterByTeam(resources, "platform", TeamFilterOptions{CaseSensitive: true})
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByTeams_ORSemantics(t *testing.T) {
	resources := []Resource{
		teamResource("r1", "platform"),
		teamResource("r2", "security"),
		teamResource("r3", "data"),
	}
	got := FilterByTeams(resources, []string{"platform", "security"}, DefaultTeamFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildTeamIndex_Lookup(t *testing.T) {
	resources := []Resource{teamResource("r1", "platform"), teamResource("r2", "security")}
	idx := BuildTeamIndex(resources)
	got := idx.Lookup("Platform")
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestBuildTeamIndex_LookupMissing(t *testing.T) {
	idx := BuildTeamIndex([]Resource{teamResource("r1", "platform")})
	if idx.Lookup("unknown") != nil {
		t.Fatal("expected nil for missing team")
	}
}
