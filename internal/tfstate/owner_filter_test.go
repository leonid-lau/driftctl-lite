package tfstate

import (
	"testing"
)

func ownerResource(id, owner string) Resource {
	return Resource{
		ID:   id,
		Type: "aws_instance",
		Metadata: map[string]string{"owner": owner},
		Attributes: map[string]interface{}{},
	}
}

func TestFilterByOwner_Match(t *testing.T) {
	resources := []Resource{
		ownerResource("r1", "team-a"),
		ownerResource("r2", "team-b"),
		ownerResource("r3", "Team-A"),
	}
	got := FilterByOwner(resources, "team-a", DefaultOwnerFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByOwner_EmptyOwner_ReturnsAll(t *testing.T) {
	resources := []Resource{ownerResource("r1", "team-a"), ownerResource("r2", "team-b")}
	got := FilterByOwner(resources, "", DefaultOwnerFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByOwner_CaseSensitive(t *testing.T) {
	resources := []Resource{ownerResource("r1", "Team-A"), ownerResource("r2", "team-a")}
	opts := OwnerFilterOptions{CaseSensitive: true}
	got := FilterByOwner(resources, "team-a", opts)
	if len(got) != 1 || got[0].ID != "r2" {
		t.Fatalf("expected only r2, got %v", got)
	}
}

func TestFilterByOwners_ORSemantics(t *testing.T) {
	resources := []Resource{
		ownerResource("r1", "team-a"),
		ownerResource("r2", "team-b"),
		ownerResource("r3", "team-c"),
	}
	got := FilterByOwners(resources, []string{"team-a", "team-c"}, DefaultOwnerFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildOwnerIndex_Lookup(t *testing.T) {
	resources := []Resource{ownerResource("r1", "team-a"), ownerResource("r2", "Team-A")}
	idx := BuildOwnerIndex(resources)
	got := idx.Lookup("team-a")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildOwnerIndex_LookupMissing(t *testing.T) {
	idx := BuildOwnerIndex([]Resource{ownerResource("r1", "team-a")})
	if got := idx.Lookup("team-z"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}
