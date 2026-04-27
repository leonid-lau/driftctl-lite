package tfstate

import (
	"testing"
)

func categoryResource(id, category string) Resource {
	attrs := map[string]interface{}{}
	if category != "" {
		attrs["category"] = category
	}
	return Resource{ID: id, Type: "test_resource", Attributes: attrs}
}

func TestFilterByCategory_Match(t *testing.T) {
	resources := []Resource{
		categoryResource("r1", "compute"),
		categoryResource("r2", "storage"),
		categoryResource("r3", "compute"),
	}
	opts := DefaultCategoryFilterOptions()
	got := FilterByCategory(resources, "compute", opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCategory_EmptyCategory_ReturnsAll(t *testing.T) {
	resources := []Resource{
		categoryResource("r1", "compute"),
		categoryResource("r2", "storage"),
	}
	opts := DefaultCategoryFilterOptions()
	got := FilterByCategory(resources, "", opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCategory_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		categoryResource("r1", "Compute"),
		categoryResource("r2", "storage"),
	}
	opts := DefaultCategoryFilterOptions()
	got := FilterByCategory(resources, "compute", opts)
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %+v", got)
	}
}

func TestFilterByCategory_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		categoryResource("r1", "Compute"),
	}
	opts := FilterOptions{CaseInsensitive: false}
	got := FilterByCategory(resources, "compute", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByCategories_ORSemantics(t *testing.T) {
	resources := []Resource{
		categoryResource("r1", "compute"),
		categoryResource("r2", "storage"),
		categoryResource("r3", "network"),
	}
	opts := DefaultCategoryFilterOptions()
	got := FilterByCategories(resources, []string{"compute", "storage"}, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCategories_EmptySlice_ReturnsAll(t *testing.T) {
	resources := []Resource{
		categoryResource("r1", "compute"),
		categoryResource("r2", "storage"),
	}
	opts := DefaultCategoryFilterOptions()
	got := FilterByCategories(resources, nil, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
