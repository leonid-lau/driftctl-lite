package tfstate

import (
	"sort"
	"testing"
)

func idxCategoryResource(id, category string) Resource {
	return Resource{
		ID:   id,
		Type: "test_resource",
		Attributes: map[string]interface{}{
			"category": category,
		},
	}
}

func TestBuildCategoryIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxCategoryResource("r1", "compute"),
		idxCategoryResource("r2", "storage"),
		idxCategoryResource("r3", "compute"),
	}
	idx := BuildCategoryIndex(resources)
	got := idx.Lookup("compute")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildCategoryIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{
		idxCategoryResource("r1", "Storage"),
	}
	idx := BuildCategoryIndex(resources)
	got := idx.Lookup("storage")
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %+v", got)
	}
}

func TestBuildCategoryIndex_LookupMissing(t *testing.T) {
	idx := BuildCategoryIndex(nil)
	got := idx.Lookup("compute")
	if got != nil {
		t.Fatalf("expected nil, got %+v", got)
	}
}

func TestBuildCategoryIndex_Categories(t *testing.T) {
	resources := []Resource{
		idxCategoryResource("r1", "compute"),
		idxCategoryResource("r2", "storage"),
		idxCategoryResource("r3", "compute"),
	}
	idx := BuildCategoryIndex(resources)
	keys := idx.Categories()
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "compute" || keys[1] != "storage" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestBuildCategoryIndex_EmptyInput(t *testing.T) {
	idx := BuildCategoryIndex([]Resource{})
	if len(idx.Categories()) != 0 {
		t.Fatal("expected empty index")
	}
}
