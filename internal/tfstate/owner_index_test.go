package tfstate

import (
	"testing"
)

func idxOwnerResource(id, owner string) Resource {
	return Resource{ID: id, Type: "aws_s3_bucket", Attributes: map[string]interface{}{"owner": owner}}
}

func TestBuildOwnerIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxOwnerResource("r1", "team-a"),
		idxOwnerResource("r2", "team-b"),
		idxOwnerResource("r3", "team-a"),
	}
	idx := BuildOwnerIndex(resources)
	got := idx.Lookup("team-a")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildOwnerIndex_LookupMissing(t *testing.T) {
	idx := BuildOwnerIndex([]Resource{})
	got := idx.Lookup("team-x")
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestBuildOwnerIndex_Owners(t *testing.T) {
	resources := []Resource{
		idxOwnerResource("r1", "team-a"),
		idxOwnerResource("r2", "team-b"),
	}
	idx := BuildOwnerIndex(resources)
	owners := idx.Owners()
	if len(owners) != 2 {
		t.Fatalf("expected 2 owners, got %d", len(owners))
	}
}
