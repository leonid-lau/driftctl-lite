package tfstate

import (
	"testing"
)

func createdResource(id, createdBy string) Resource {
	attrs := map[string]interface{}{}
	if createdBy != "" {
		attrs["created_by"] = createdBy
	}
	return Resource{ID: id, Type: "aws_instance", Attributes: attrs}
}

func TestFilterByCreated_Match(t *testing.T) {
	resources := []Resource{
		createdResource("r1", "terraform"),
		createdResource("r2", "pulumi"),
		createdResource("r3", "terraform"),
	}
	got := FilterByCreated(resources, "terraform", DefaultCreatedFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCreated_EmptyCreatedBy_ReturnsAll(t *testing.T) {
	resources := []Resource{
		createdResource("r1", "terraform"),
		createdResource("r2", "pulumi"),
	}
	got := FilterByCreated(resources, "", DefaultCreatedFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCreated_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		createdResource("r1", "Terraform"),
	}
	got := FilterByCreated(resources, "terraform", DefaultCreatedFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByCreated_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		createdResource("r1", "Terraform"),
	}
	opts := CreatedFilterOptions{CaseSensitive: true}
	got := FilterByCreated(resources, "terraform", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByCreatedValues_ORSemantics(t *testing.T) {
	resources := []Resource{
		createdResource("r1", "terraform"),
		createdResource("r2", "pulumi"),
		createdResource("r3", "ansible"),
	}
	got := FilterByCreatedValues(resources, []string{"terraform", "pulumi"}, DefaultCreatedFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCreated_MissingAttribute(t *testing.T) {
	resources := []Resource{
		createdResource("r1", ""),
		createdResource("r2", "terraform"),
	}
	got := FilterByCreated(resources, "terraform", DefaultCreatedFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
	if got[0].ID != "r2" {
		t.Errorf("expected r2, got %s", got[0].ID)
	}
}
