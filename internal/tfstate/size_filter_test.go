package tfstate

import (
	"testing"
)

func sizeResource(id, size string) Resource {
	return Resource{
		Type: "aws_instance",
		Name: id,
		Attributes: map[string]interface{}{"size": size},
	}
}

func TestFilterBySize_Match(t *testing.T) {
	resources := []Resource{
		sizeResource("a", "small"),
		sizeResource("b", "medium"),
		sizeResource("c", "large"),
	}
	got := FilterBySize(resources, "medium", DefaultSizeFilterOptions())
	if len(got) != 1 || got[0].Name != "b" {
		t.Fatalf("expected [b], got %v", got)
	}
}

func TestFilterBySize_EmptySize_ReturnsAll(t *testing.T) {
	resources := []Resource{sizeResource("a", "small"), sizeResource("b", "large")}
	got := FilterBySize(resources, "", DefaultSizeFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterBySize_CaseInsensitive(t *testing.T) {
	resources := []Resource{sizeResource("a", "Large")}
	got := FilterBySize(resources, "large", DefaultSizeFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterBySize_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{sizeResource("a", "Large")}
	opts := SizeFilterOptions{CaseSensitive: true}
	got := FilterBySize(resources, "large", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterBySizes_ORSemantics(t *testing.T) {
	resources := []Resource{
		sizeResource("a", "small"),
		sizeResource("b", "medium"),
		sizeResource("c", "large"),
	}
	got := FilterBySizes(resources, []string{"small", "large"}, DefaultSizeFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildSizeIndex_Lookup(t *testing.T) {
	resources := []Resource{
		sizeResource("a", "small"),
		sizeResource("b", "SMALL"),
		sizeResource("c", "large"),
	}
	idx := BuildSizeIndex(resources)
	got := idx.Lookup("small")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildSizeIndex_LookupMissing(t *testing.T) {
	idx := BuildSizeIndex([]Resource{sizeResource("a", "small")})
	if idx.Lookup("xlarge") != nil {
		t.Fatal("expected nil for missing size")
	}
}

func TestBuildSizeIndex_Sizes(t *testing.T) {
	resources := []Resource{
		sizeResource("a", "small"),
		sizeResource("b", "large"),
	}
	idx := BuildSizeIndex(resources)
	sizes := idx.Sizes()
	if len(sizes) != 2 {
		t.Fatalf("expected 2 sizes, got %d", len(sizes))
	}
}
