package tfstate

import (
	"testing"
)

func idxWSResource(workspace string) Resource {
	return Resource{
		Type: "aws_s3_bucket",
		Attributes: map[string]interface{}{"workspace": workspace},
	}
}

func TestBuildWorkspaceIndex_Lookup(t *testing.T) {
	resources := []Resource{idxWSResource("prod"), idxWSResource("dev")}
	idx := BuildWorkspaceIndex(resources)
	got := idx.Lookup("prod")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildWorkspaceIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxWSResource("Prod")}
	idx := BuildWorkspaceIndex(resources)
	got := idx.Lookup("prod")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildWorkspaceIndex_LookupMissing(t *testing.T) {
	idx := BuildWorkspaceIndex([]Resource{idxWSResource("dev")})
	got := idx.Lookup("prod")
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestBuildWorkspaceIndex_Workspaces(t *testing.T) {
	resources := []Resource{idxWSResource("prod"), idxWSResource("dev"), idxWSResource("prod")}
	idx := BuildWorkspaceIndex(resources)
	ws := idx.Workspaces()
	if len(ws) != 2 {
		t.Fatalf("expected 2 unique workspaces, got %d", len(ws))
	}
}

func TestBuildWorkspaceIndex_EmptyInput(t *testing.T) {
	idx := BuildWorkspaceIndex(nil)
	if len(idx.Workspaces()) != 0 {
		t.Fatal("expected empty index")
	}
}
