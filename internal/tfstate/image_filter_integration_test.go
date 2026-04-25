package tfstate

import "testing"

// TestImageFilterAndIndex_RoundTrip verifies that BuildImageIndex and
// FilterByImage produce consistent results for the same resource set.
func TestImageFilterAndIndex_RoundTrip(t *testing.T) {
	resources := []Resource{
		idxImageResource("ami-prod-001"),
		idxImageResource("ami-prod-002"),
		idxImageResource("ami-staging-001"),
	}

	target := "ami-prod-001"

	// Filter approach
	filtered := FilterByImage(resources, target, DefaultImageFilterOptions())

	// Index approach
	idx := BuildImageIndex(resources)
	indexed := idx.Lookup(target)

	if len(filtered) != len(indexed) {
		t.Fatalf("filter returned %d but index returned %d", len(filtered), len(indexed))
	}

	if len(filtered) != 1 {
		t.Fatalf("expected exactly 1 match, got %d", len(filtered))
	}

	v, _ := filtered[0].Attributes["image"].(string)
	if v != target {
		t.Fatalf("unexpected image value: %s", v)
	}
}
