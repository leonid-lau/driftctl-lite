package tfstate

import "testing"

func instanceResource(id, instance string) Resource {
	attrs := map[string]interface{}{}
	if instance != "" {
		attrs["instance"] = instance
	}
	return Resource{ID: id, Type: "aws_instance", Attributes: attrs}
}

func TestFilterByInstance_Match(t *testing.T) {
	resources := []Resource{
		instanceResource("r1", "web"),
		instanceResource("r2", "db"),
		instanceResource("r3", "cache"),
	}
	got := FilterByInstance(resources, "db", DefaultInstanceFilterOptions())
	if len(got) != 1 || got[0].ID != "r2" {
		t.Fatalf("expected [r2], got %v", got)
	}
}

func TestFilterByInstance_EmptyInstance_ReturnsAll(t *testing.T) {
	resources := []Resource{
		instanceResource("r1", "web"),
		instanceResource("r2", "db"),
	}
	got := FilterByInstance(resources, "", DefaultInstanceFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByInstance_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		instanceResource("r1", "Web"),
		instanceResource("r2", "db"),
	}
	got := FilterByInstance(resources, "web", DefaultInstanceFilterOptions())
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected [r1], got %v", got)
	}
}

func TestFilterByInstance_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		instanceResource("r1", "Web"),
	}
	opts := InstanceFilterOptions{CaseSensitive: true}
	got := FilterByInstance(resources, "web", opts)
	if len(got) != 0 {
		t.Fatalf("expected no results, got %v", got)
	}
}

func TestFilterByInstances_ORSemantics(t *testing.T) {
	resources := []Resource{
		instanceResource("r1", "web"),
		instanceResource("r2", "db"),
		instanceResource("r3", "cache"),
	}
	got := FilterByInstances(resources, []string{"web", "cache"}, DefaultInstanceFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestFilterByInstances_EmptySlice_ReturnsAll(t *testing.T) {
	resources := []Resource{
		instanceResource("r1", "web"),
		instanceResource("r2", "db"),
	}
	got := FilterByInstances(resources, []string{}, DefaultInstanceFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByInstance_MissingAttribute(t *testing.T) {
	resources := []Resource{
		instanceResource("r1", ""), // no instance attribute set
	}
	got := FilterByInstance(resources, "web", DefaultInstanceFilterOptions())
	if len(got) != 0 {
		t.Fatalf("expected no results, got %v", got)
	}
}
