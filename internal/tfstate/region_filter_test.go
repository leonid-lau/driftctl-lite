package tfstate

import (
	"testing"
)

func regionResource(id, region string) Resource {
	return Resource{ID: id, Type: "aws_instance", Region: region}
}

func TestFilterByRegion_Match(t *testing.T) {
	resources := []Resource{
		regionResource("a", "us-east-1"),
		regionResource("b", "eu-west-1"),
		regionResource("c", "us-east-1"),
	}
	got := FilterByRegion(resources, "us-east-1", DefaultRegionFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByRegion_EmptyRegion_ReturnsAll(t *testing.T) {
	resources := []Resource{
		regionResource("a", "us-east-1"),
		regionResource("b", "eu-west-1"),
	}
	got := FilterByRegion(resources, "", DefaultRegionFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByRegion_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		regionResource("a", "US-EAST-1"),
	}
	got := FilterByRegion(resources, "us-east-1", DefaultRegionFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByRegion_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		regionResource("a", "US-EAST-1"),
	}
	opts := RegionFilterOptions{CaseInsensitive: false}
	got := FilterByRegion(resources, "us-east-1", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByRegions_ORSemantics(t *testing.T) {
	resources := []Resource{
		regionResource("a", "us-east-1"),
		regionResource("b", "eu-west-1"),
		regionResource("c", "ap-southeast-1"),
	}
	got := FilterByRegions(resources, []string{"us-east-1", "eu-west-1"}, DefaultRegionFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByRegions_EmptySlice_ReturnsAll(t *testing.T) {
	resources := []Resource{
		regionResource("a", "us-east-1"),
		regionResource("b", "eu-west-1"),
	}
	got := FilterByRegions(resources, nil, DefaultRegionFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
