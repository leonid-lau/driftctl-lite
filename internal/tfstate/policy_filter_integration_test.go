package tfstate

import (
	"testing"
)

// TestPolicyFilterAndIndex_RoundTrip verifies that FilterByPolicy and BuildPolicyIndex
// return consistent results for the same inputs.
func TestPolicyFilterAndIndex_RoundTrip(t *testing.T) {
	resources := []Resource{
		{ID: "r1", Type: "aws_iam_policy", Attributes: map[string]string{"policy": "ReadOnly"}},
		{ID: "r2", Type: "aws_iam_policy", Attributes: map[string]string{"policy": "FullAccess"}},
		{ID: "r3", Type: "aws_iam_policy", Attributes: map[string]string{"policy": "ReadOnly"}},
	}

	filtered := FilterByPolicy(resources, "ReadOnly", DefaultPolicyFilterOptions())
	if len(filtered) != 2 {
		t.Fatalf("filter: expected 2, got %d", len(filtered))
	}

	idx := BuildPolicyIndex(resources)
	looked := idx.Lookup("ReadOnly")
	if len(looked) != 2 {
		t.Fatalf("index: expected 2, got %d", len(looked))
	}

	for i, r := range filtered {
		if r.ID != looked[i].ID {
			t.Errorf("mismatch at %d: filter=%s index=%s", i, r.ID, looked[i].ID)
		}
	}
}
