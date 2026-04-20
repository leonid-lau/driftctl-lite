package tfstate

import (
	"testing"
)

func depResource(dep string) Resource {
	attrs := map[string]interface{}{}
	if dep != "" {
		attrs["dependency"] = dep
	}
	return Resource{Type: "aws_instance", Attributes: attrs}
}

func TestFilterByDependency_Match(t *testing.T) {
	resources := []Resource{depResource("moduleA"), depResource("moduleB"), depResource("moduleA")}
	got := FilterByDependency(resources, "moduleA", DefaultDependencyFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByDependency_EmptyDep_ReturnsAll(t *testing.T) {
	resources := []Resource{depResource("moduleA"), depResource("moduleB")}
	got := FilterByDependency(resources, "", DefaultDependencyFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByDependency_CaseInsensitive(t *testing.T) {
	resources := []Resource{depResource("ModuleA")}
	got := FilterByDependency(resources, "modulea", DefaultDependencyFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByDependency_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{depResource("ModuleA")}
	opts := FilterOptions{CaseSensitive: true}
	got := FilterByDependency(resources, "modulea", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByDependencies_ORSemantics(t *testing.T) {
	resources := []Resource{depResource("moduleA"), depResource("moduleB"), depResource("moduleC")}
	got := FilterByDependencies(resources, []string{"moduleA", "moduleB"}, DefaultDependencyFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByDependencies_EmptyList_ReturnsAll(t *testing.T) {
	resources := []Resource{depResource("moduleA"), depResource("moduleB")}
	got := FilterByDependencies(resources, nil, DefaultDependencyFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
