package tfstate

import (
	"sort"
	"testing"
)

func idxProjectResource(id, project string) Resource {
	return Resource{
		ID:   id,
		Type: "aws_s3_bucket",
		Attributes: map[string]string{
			"project": project,
		},
	}
}

func TestBuildProjectIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxProjectResource("r1", "alpha"),
		idxProjectResource("r2", "beta"),
		idxProjectResource("r3", "alpha"),
	}
	idx := BuildProjectIndex(resources)
	got := idx.Lookup("alpha")
	if len(got) != 2 {
		t.Fatalf("expected 2 resources for project alpha, got %d", len(got))
	}
}

func TestBuildProjectIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{
		idxProjectResource("r1", "Alpha"),
	}
	idx := BuildProjectIndex(resources)
	got := idx.Lookup("ALPHA")
	if len(got) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got))
	}
}

func TestBuildProjectIndex_LookupMissing(t *testing.T) {
	idx := BuildProjectIndex([]Resource{})
	if got := idx.Lookup("nonexistent"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuildProjectIndex_Projects(t *testing.T) {
	resources := []Resource{
		idxProjectResource("r1", "alpha"),
		idxProjectResource("r2", "beta"),
		idxProjectResource("r3", "gamma"),
	}
	idx := BuildProjectIndex(resources)
	projects := idx.Projects()
	sort.Strings(projects)
	if len(projects) != 3 {
		t.Fatalf("expected 3 projects, got %d: %v", len(projects), projects)
	}
}

func TestBuildProjectIndex_EmptyProject_Skipped(t *testing.T) {
	resources := []Resource{
		{ID: "r1", Type: "aws_s3_bucket", Attributes: map[string]string{}},
		idxProjectResource("r2", "alpha"),
	}
	idx := BuildProjectIndex(resources)
	if len(idx.Projects()) != 1 {
		t.Fatalf("expected 1 project, got %d", len(idx.Projects()))
	}
}
