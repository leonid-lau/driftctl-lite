package drift_test

import (
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/tfstate"
)

func stateIndex(resources []tfstate.Resource) map[string]tfstate.Resource {
	idx := make(map[string]tfstate.Resource, len(resources))
	for _, r := range resources {
		id, _ := r.Attributes["id"].(string)
		if id == "" {
			id = r.Name
		}
		idx[id] = r
	}
	return idx
}

func TestDetect_NoDrift(t *testing.T) {
	det := drift.NewDetector()  
	si := map[string]tfstate.Resource{
		"i-123": {Name: "web", Type: "aws_instance", Attributes: map[string]interface{}{"id": "i-123", "ami": "ami-abc"}},
	}
	li := map[string]map[string]interface{}{
		"i-123": {"id": "i-123", "ami": "ami-abc"},
	}
	results := det.Detect(si, li)
	if len(results) != 0 {
		t.Fatalf("expected no drift, got %d results: %v", len(results), results)
	}
}

func TestDetect_MissingResource(t *testing.T) {
	det := drift.NewDetector()
	si := map[string]tfstate.Resource{
		"i-999": {Name: "db", Type: "aws_instance", Attributes: map[string]interface{}{"id": "i-999"}},
	}
	li := map[string]map[string]interface{}{}
	results := det.Detect(si, li)
	if len(results) != 1 || results[0].Type != drift.DriftMissing {
		t.Fatalf("expected one MISSING drift, got %v", results)
	}
}

func TestDetect_ExtraResource(t *testing.T) {
	det := drift.NewDetector()
	si := map[string]tfstate.Resource{}
	li := map[string]map[string]interface{}{
		"i-extra": {"id": "i-extra"},
	}
	results := det.Detect(si, li)
	if len(results) != 1 || results[0].Type != drift.DriftExtra {
		t.Fatalf("expected one EXTRA drift, got %v", results)
	}
}

func TestDetect_ModifiedAttribute(t *testing.T) {
	det := drift.NewDetector()
	si := map[string]tfstate.Resource{
		"i-001": {Name: "app", Type: "aws_instance", Attributes: map[string]interface{}{"id": "i-001", "instance_type": "t2.micro"}},
	}
	li := map[string]map[string]interface{}{
		"i-001": {"id": "i-001", "instance_type": "t3.medium"},
	}
	results := det.Detect(si, li)
	if len(results) != 1 || results[0].Type != drift.DriftModified {
		t.Fatalf("expected one MODIFIED drift, got %v", results)
	}
	if results[0].Attribute != "instance_type" {
		t.Errorf("expected attribute 'instance_type', got %q", results[0].Attribute)
	}
}

// TestDetect_MultipleModifiedAttributes verifies that each differing attribute
// produces a separate DriftModified result when multiple attributes diverge.
func TestDetect_MultipleModifiedAttributes(t *testing.T) {
	det := drift.NewDetector()
	si := map[string]tfstate.Resource{
		"i-002": {Name: "app2", Type: "aws_instance", Attributes: map[string]interface{}{
			"id":            "i-002",
			"instance_type": "t2.micro",
			"ami":           "ami-old",
		}},
	}
	li := map[string]map[string]interface{}{
		"i-002": {"id": "i-002", "instance_type": "t3.large", "ami": "ami-new"},
	}
	results := det.Detect(si, li)
	modifiedCount := 0
	for _, r := range results {
		if r.Type == drift.DriftModified {
			modifiedCount++
		}
	}
	if modifiedCount != 2 {
		t.Fatalf("expected 2 MODIFIED drift results, got %d: %v", modifiedCount, results)
	}
}

func TestDriftResult_String(t *testing.T) {
	cases := []struct {
		result   drift.DriftResult
		contains string
	}{
		{drift.DriftResult{Type: drift.DriftMissing, ResourceID: "r1"}, "MISSING"},
		{drift.DriftResult{Type: drift.DriftExtra, ResourceID: "r2"}, "EXTRA"},
		{drift.DriftResult{Type: drift.DriftModified, ResourceID: "r3", Attribute: "ami"}, "MODIFIED"},
	}
	for _, tc := range cases {
		s := tc.result.String()
		if s == "" {
			t.Errorf("expected non-empty string for %v", tc.result)
		}
	}
}
