package tfstate

import (
	"testing"
)

func idxRegionResource(region string) Resource {
	return Resource{
		Type:       "aws_vpc",
		Attributes: map[string]string{"region": region},
	}
}

func TestBuildRegionIndex_Lookup(t *testing.T) {
	resources := []Resource{idxRegionResource("us-east-1"), idxRegionResource("eu-west-1")}
	idx := BuildRegionIndex(resources)
	got := idx.Lookup("us-east-1")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildRegionIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxRegionResource("US-EAST-1")}
	idx := BuildRegionIndex(resources)
	if got := idx.Lookup("us-east-1"); len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildRegionIndex_LookupMissing(t *testing.T) {
	idx := BuildRegionIndex([]Resource{idxRegionResource("us-east-1")})
	if got := idx.Lookup("ap-south-1"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuildRegionIndex_Regions(t *testing.T) {
	resources := []Resource{idxRegionResource("us-east-1"), idxRegionResource("eu-west-1")}
	idx := BuildRegionIndex(resources)
	if len(idx.Regions()) != 2 {
		t.Fatalf("expected 2 regions, got %d", len(idx.Regions()))
	}
}

func TestBuildRegionIndex_EmptyInput(t *testing.T) {
	idx := BuildRegionIndex(nil)
	if len(idx.Regions()) != 0 {
		t.Fatal("expected empty")
	}
}
