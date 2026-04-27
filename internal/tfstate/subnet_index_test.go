package tfstate

import (
	"testing"
)

func idxSubnetResource(subnet string) Resource {
	return Resource{
		Type: "aws_subnet",
		Attributes: map[string]interface{}{
			"subnet": subnet,
		},
	}
}

func idxSubnetIDResource(subnetID string) Resource {
	return Resource{
		Type: "aws_subnet",
		Attributes: map[string]interface{}{
			"subnet_id": subnetID,
		},
	}
}

func TestBuildSubnetIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxSubnetResource("subnet-abc"),
		idxSubnetResource("subnet-def"),
	}
	idx := BuildSubnetIndex(resources)
	got := idx.Lookup("subnet-abc")
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestBuildSubnetIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{
		idxSubnetResource("Subnet-XYZ"),
	}
	idx := BuildSubnetIndex(resources)
	got := idx.Lookup("subnet-xyz")
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestBuildSubnetIndex_FallbackToSubnetID(t *testing.T) {
	resources := []Resource{
		idxSubnetIDResource("subnet-fallback"),
	}
	idx := BuildSubnetIndex(resources)
	got := idx.Lookup("subnet-fallback")
	if len(got) != 1 {
		t.Fatalf("expected 1 result via fallback, got %d", len(got))
	}
}

func TestBuildSubnetIndex_LookupMissing(t *testing.T) {
	resources := []Resource{
		idxSubnetResource("subnet-abc"),
	}
	idx := BuildSubnetIndex(resources)
	got := idx.Lookup("subnet-zzz")
	if len(got) != 0 {
		t.Fatalf("expected 0 results, got %d", len(got))
	}
}

func TestBuildSubnetIndex_Subnets(t *testing.T) {
	resources := []Resource{
		idxSubnetResource("subnet-a"),
		idxSubnetResource("subnet-b"),
		idxSubnetIDResource("subnet-c"),
	}
	idx := BuildSubnetIndex(resources)
	subnets := idx.Subnets()
	if len(subnets) != 3 {
		t.Fatalf("expected 3 subnets, got %d", len(subnets))
	}
}

func TestBuildSubnetIndex_EmptyInput(t *testing.T) {
	idx := BuildSubnetIndex(nil)
	got := idx.Lookup("subnet-x")
	if len(got) != 0 {
		t.Fatalf("expected 0 results for empty index, got %d", len(got))
	}
}
