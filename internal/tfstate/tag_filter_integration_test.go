package tfstate

import (
	"testing"
)

// TestTagFilterAndIndex_RoundTrip verifies that BuildTagIndex and FilterByTag
// produce consistent results for the same resource set.
func TestTagFilterAndIndex_RoundTrip(t *testing.T) {
	resources := []Resource{
		{ID: "r1", Type: "aws_instance", Attributes: map[string]string{
			"tags.env": "prod", "tags.team": "platform",
		}},
		{ID: "r2", Type: "aws_s3_bucket", Attributes: map[string]string{
			"tags.env": "dev",
		}},
		{ID: "r3", Type: "aws_lambda_function", Attributes: map[string]string{
			"tags.team": "platform",
		}},
	}

	filtered := FilterByTag(resources, "team", "platform")
	idx := BuildTagIndex(resources)
	indexed := idx.Lookup("team", "platform")

	if len(filtered) != len(indexed) {
		t.Fatalf("filter returned %d, index returned %d", len(filtered), len(indexed))
	}

	filterIDs := map[string]bool{}
	for _, r := range filtered {
		filterIDs[r.ID] = true
	}
	for _, r := range indexed {
		if !filterIDs[r.ID] {
			t.Errorf("index returned %q not in filter results", r.ID)
		}
	}
}
