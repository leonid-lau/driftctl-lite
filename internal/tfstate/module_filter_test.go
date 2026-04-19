package tfstate

import (
	"testing"
)

func moduleResource(module string) Resource {
	return Resource{
		Type: "aws_instance",
		Attributes: map[string]interface{}{"module": module},
	}
}

func TestFilterByModule_Match(t *testing.T) {
	resources := []Resource{moduleResource("vpc"), moduleResource("ecs"), moduleResource("vpc")}
	got := FilterByModule(resources, "vpc", DefaultModuleFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByModule_EmptyModule_ReturnsAll(t *testing.T) {
	resources := []Resource{moduleResource("vpc"), moduleResource("ecs")}
	got := FilterByModule(resources, "", DefaultModuleFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByModule_CaseInsensitive(t *testing.T) {
	resources := []Resource{moduleResource("VPC"), moduleResource("ecs")}
	got := FilterByModule(resources, "vpc", DefaultModuleFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByModule_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{moduleResource("VPC")}
	opts := ModuleFilterOptions{CaseSensitive: true}
	got := FilterByModule(resources, "vpc", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByModules_ORSemantics(t *testing.T) {
	resources := []Resource{moduleResource("vpc"), moduleResource("ecs"), moduleResource("rds")}
	got := FilterByModules(resources, []string{"vpc", "rds"}, DefaultModuleFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByModules_Empty_ReturnsAll(t *testing.T) {
	resources := []Resource{moduleResource("vpc"), moduleResource("ecs")}
	got := FilterByModules(resources, nil, DefaultModuleFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
