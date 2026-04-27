package tfstate_test

import (
	"testing"

	"github.com/your-org/driftctl-lite/internal/tfstate"
)

// TestPortFilterAndIndex_RoundTrip verifies that BuildPortIndex and FilterByPort
// produce consistent results when used together on the same resource set.
func TestPortFilterAndIndex_RoundTrip(t *testing.T) {
	resources := []tfstate.Resource{
		{
			Type: "aws_security_group_rule",
			ID:   "sgr-443",
			Attributes: map[string]interface{}{
				"port": "443",
			},
		},
		{
			Type: "aws_security_group_rule",
			ID:   "sgr-80",
			Attributes: map[string]interface{}{
				"port": "80",
			},
		},
		{
			Type: "aws_lb_listener",
			ID:   "lb-8080",
			Attributes: map[string]interface{}{
				"port": "8080",
			},
		},
	}

	targetPort := "443"

	// Filter using FilterByPort
	filtered := tfstate.FilterByPort(resources, targetPort, tfstate.DefaultPortFilterOptions())

	// Lookup using BuildPortIndex
	idx := tfstate.BuildPortIndex(resources)
	indexed := idx.Lookup(targetPort)

	if len(filtered) != len(indexed) {
		t.Fatalf("filter returned %d resources, index returned %d", len(filtered), len(indexed))
	}

	for i, r := range filtered {
		if r.ID != indexed[i].ID {
			t.Errorf("mismatch at position %d: filter=%s, index=%s", i, r.ID, indexed[i].ID)
		}
	}
}
