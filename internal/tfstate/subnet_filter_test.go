package tfstate

import (
	"testing"
)

func subnetResource(subnet string) Resource {
	return Resource{
		Type: "aws_subnet",
		ID:   subnet,
		Attributes: map[string]string{
			"subnet": subnet,
		},
	}
}

func subnetIDResource(subnetID string) Resource {
	return Resource{
		Type: "aws_instance",
		ID:   subnetID,
		Attributes: map[string]string{
			"subnet_id": subnetID,
		},
	}
}

func TestFilterBySubnet_Match(t *testing.T) {
	resources := []Resource{subnetResource("subnet-abc"), subnetResource("subnet-xyz")}
	got := FilterBySubnet(resources, "subnet-abc", DefaultSubnetFilterOptions())
	if len(got) != 1 || got[0].ID != "subnet-abc" {
		t.Fatalf("expected 1 match, got %v", got)
	}
}

func TestFilterBySubnet_EmptySubnet_ReturnsAll(t *testing.T) {
	resources := []Resource{subnetResource("subnet-a"), subnetResource("subnet-b")}
	got := FilterBySubnet(resources, "", DefaultSubnetFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterBySubnet_CaseInsensitive(t *testing.T) {
	resources := []Resource{subnetResource("Subnet-ABC")}
	got := FilterBySubnet(resources, "subnet-abc", DefaultSubnetFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1 match, got %d", len(got))
	}
}

func TestFilterBySubnet_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{subnetResource("Subnet-ABC")}
	opts := SubnetFilterOptions{CaseSensitive: true}
	got := FilterBySubnet(resources, "subnet-abc", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(got))
	}
}

func TestFilterBySubnets_ORSemantics(t *testing.T) {
	resources := []Resource{
		subnetResource("subnet-a"),
		subnetResource("subnet-b"),
		subnetResource("subnet-c"),
	}
	got := FilterBySubnets(resources, []string{"subnet-a", "subnet-c"}, DefaultSubnetFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildSubnetIndex_Lookup(t *testing.T) {
	resources := []Resource{subnetResource("subnet-1"), subnetIDResource("subnet-2")}
	idx := BuildSubnetIndex(resources)
	if res := idx.Lookup("subnet-1"); len(res) != 1 {
		t.Fatalf("expected 1 result for subnet-1, got %d", len(res))
	}
	if res := idx.Lookup("subnet-2"); len(res) != 1 {
		t.Fatalf("expected 1 result for subnet-2, got %d", len(res))
	}
}

func TestBuildSubnetIndex_LookupMissing(t *testing.T) {
	idx := BuildSubnetIndex([]Resource{subnetResource("subnet-x")})
	if res := idx.Lookup("subnet-z"); len(res) != 0 {
		t.Fatalf("expected 0, got %d", len(res))
	}
}

func TestBuildSubnetIndex_Subnets(t *testing.T) {
	resources := []Resource{subnetResource("subnet-a"), subnetResource("subnet-b")}
	idx := BuildSubnetIndex(resources)
	if len(idx.Subnets()) != 2 {
		t.Fatalf("expected 2 subnets, got %d", len(idx.Subnets()))
	}
}
