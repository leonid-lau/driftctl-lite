package tfstate

import (
	"testing"
)

func wsResource(workspace string) Resource {
	return Resource{
		Type: "aws_instance",
		Attributes: map[string]interface{}{"workspace": workspace},
	}
}

func TestFilterByWorkspace_Match(t *testing.T) {
	resources := []Resource{wsResource("prod"), wsResource("staging"), wsResource("dev")}
	got := FilterByWorkspace(resources, "prod", DefaultWorkspaceFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByWorkspace_EmptyWorkspace_ReturnsAll(t *testing.T) {
	resources := []Resource{wsResource("prod"), wsResource("dev")}
	got := FilterByWorkspace(resources, "", DefaultWorkspaceFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByWorkspace_CaseInsensitive(t *testing.T) {
	resources := []Resource{wsResource("Prod"), wsResource("dev")}
	got := FilterByWorkspace(resources, "prod", DefaultWorkspaceFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByWorkspace_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{wsResource("Prod")}
	opts := WorkspaceFilterOptions{CaseSensitive: true}
	got := FilterByWorkspace(resources, "prod", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByWorkspaces_ORSemantics(t *testing.T) {
	resources := []Resource{wsResource("prod"), wsResource("staging"), wsResource("dev")}
	got := FilterByWorkspaces(resources, []string{"prod", "dev"}, DefaultWorkspaceFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByWorkspaces_Empty_ReturnsAll(t *testing.T) {
	resources := []Resource{wsResource("prod"), wsResource("dev")}
	got := FilterByWorkspaces(resources, nil, DefaultWorkspaceFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
