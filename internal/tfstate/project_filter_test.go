package tfstate

import (
	"testing"
)

func projectResource(project string) Resource {
	return Resource{
		Type: "aws_project",
		Name: "example",
		Attributes: map[string]interface{}{"project": project},
	}
}

func TestFilterByProject_Match(t *testing.T) {
	resources := []Resource{
		projectResource("alpha"),
		projectResource("beta"),
		projectResource("gamma"),
	}
	got := FilterByProject(resources, "beta", DefaultProjectFilterOptions())
	if len(got) != 1 || got[0].Attributes["project"] != "beta" {
		t.Fatalf("expected 1 resource with project=beta, got %v", got)
	}
}

func TestFilterByProject_EmptyProject_ReturnsAll(t *testing.T) {
	resources := []Resource{
		projectResource("alpha"),
		projectResource("beta"),
	}
	got := FilterByProject(resources, "", DefaultProjectFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByProject_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		projectResource("Alpha"),
		projectResource("beta"),
	}
	got := FilterByProject(resources, "alpha", DefaultProjectFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got))
	}
}

func TestFilterByProject_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		projectResource("Alpha"),
	}
	opts := ProjectFilterOptions{CaseSensitive: true}
	got := FilterByProject(resources, "alpha", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0 resources, got %d", len(got))
	}
}

func TestFilterByProjects_ORSemantics(t *testing.T) {
	resources := []Resource{
		projectResource("alpha"),
		projectResource("beta"),
		projectResource("gamma"),
	}
	got := FilterByProjects(resources, []string{"alpha", "gamma"}, DefaultProjectFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByProjects_EmptyList_ReturnsAll(t *testing.T) {
	resources := []Resource{
		projectResource("alpha"),
		projectResource("beta"),
	}
	got := FilterByProjects(resources, nil, DefaultProjectFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}
