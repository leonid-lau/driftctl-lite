package tfstate

import (
	"testing"
)

func componentResource(component string) Resource {
	attrs := map[string]interface{}{}
	if component != "" {
		attrs["component"] = component
	}
	return Resource{Type: "aws_instance", ID: component, Attributes: attrs}
}

func TestFilterByComponent_Match(t *testing.T) {
	resources := []Resource{
		componentResource("frontend"),
		componentResource("backend"),
		componentResource("database"),
	}
	got := FilterByComponent(resources, "backend", DefaultComponentFilterOptions())
	if len(got) != 1 || got[0].ID != "backend" {
		t.Fatalf("expected 1 backend resource, got %v", got)
	}
}

func TestFilterByComponent_EmptyComponent_ReturnsAll(t *testing.T) {
	resources := []Resource{componentResource("frontend"), componentResource("backend")}
	got := FilterByComponent(resources, "", DefaultComponentFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByComponent_CaseInsensitive(t *testing.T) {
	resources := []Resource{componentResource("Frontend")}
	got := FilterByComponent(resources, "frontend", DefaultComponentFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByComponent_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{componentResource("Frontend")}
	opts := ComponentFilterOptions{CaseSensitive: true}
	got := FilterByComponent(resources, "frontend", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByComponents_ORSemantics(t *testing.T) {
	resources := []Resource{
		componentResource("frontend"),
		componentResource("backend"),
		componentResource("worker"),
	}
	got := FilterByComponents(resources, []string{"frontend", "worker"}, DefaultComponentFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildComponentIndex_Lookup(t *testing.T) {
	resources := []Resource{componentResource("api"), componentResource("api"), componentResource("ui")}
	idx := BuildComponentIndex(resources)
	if len(idx.Lookup("api")) != 2 {
		t.Fatalf("expected 2 api entries")
	}
	if len(idx.Lookup("ui")) != 1 {
		t.Fatalf("expected 1 ui entry")
	}
	if len(idx.Lookup("missing")) != 0 {
		t.Fatalf("expected 0 for missing")
	}
}

func TestBuildComponentIndex_Components(t *testing.T) {
	resources := []Resource{componentResource("alpha"), componentResource("beta")}
	idx := BuildComponentIndex(resources)
	comps := idx.Components()
	if len(comps) != 2 {
		t.Fatalf("expected 2 components, got %d", len(comps))
	}
}
