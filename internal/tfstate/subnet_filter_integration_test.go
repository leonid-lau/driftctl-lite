package tfstate

import (
	"testing"
)

func TestSubnetFilterAndIndex_RoundTrip(t *testing.T) {
	resources := []Resource{
		subnetResource("subnet-aaa"),
		subnetResource("subnet-bbb"),
		subnetIDResource("subnet-ccc"),
		{
			Type:       "aws_instance",
			Attributes: map[string]interface{}{"subnet": "subnet-aaa"},
		},
	}

	// Filter to a specific subnet via FilterBySubnet
	filtered := FilterBySubnet(resources, "subnet-aaa", DefaultSubnetFilterOptions())
	if len(filtered) != 2 {
		t.Fatalf("expected 2 resources matching subnet-aaa, got %d", len(filtered))
	}

	// Build index and verify lookup returns same count
	idx := BuildSubnetIndex(resources)
	looked := idx.Lookup("subnet-aaa")
	if len(looked) != len(filtered) {
		t.Fatalf("index lookup returned %d, filter returned %d; expected match", len(looked), len(filtered))
	}

	// Verify fallback via subnet_id
	fallbackFiltered := FilterBySubnet(resources, "subnet-ccc", DefaultSubnetFilterOptions())
	fallbackLooked := idx.Lookup("subnet-ccc")
	if len(fallbackFiltered) != 1 {
		t.Fatalf("expected 1 fallback filter result, got %d", len(fallbackFiltered))
	}
	if len(fallbackLooked) != 1 {
		t.Fatalf("expected 1 fallback index result, got %d", len(fallbackLooked))
	}

	// Verify multi-subnet filter
	multi := FilterBySubnets(resources, []string{"subnet-aaa", "subnet-bbb"}, DefaultSubnetFilterOptions())
	if len(multi) != 3 {
		t.Fatalf("expected 3 resources for multi-subnet filter, got %d", len(multi))
	}
}
