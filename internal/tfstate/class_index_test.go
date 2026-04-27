package tfstate

import (
	"testing"
)

func idxClassResource(id, classVal string) Resource {
	return Resource{
		ID:   id,
		Type: "aws_instance",
		Attributes: map[string]interface{}{
			"class": classVal,
		},
	}
}

func TestBuildClassIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxClassResource("r1", "standard"),
		idxClassResource("r2", "premium"),
		idxClassResource("r3", "standard"),
	}

	idx := BuildClassIndex(resources)

	got := idx.Lookup("standard")
	if len(got) != 2 {
		t.Fatalf("expected 2 resources for class 'standard', got %d", len(got))
	}
}

func TestBuildClassIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{
		idxClassResource("r1", "Premium"),
	}

	idx := BuildClassIndex(resources)

	got := idx.Lookup("premium")
	if len(got) != 1 {
		t.Fatalf("expected 1 resource for class 'premium' (case-insensitive), got %d", len(got))
	}
	if got[0].ID != "r1" {
		t.Errorf("expected resource ID 'r1', got %q", got[0].ID)
	}
}

func TestBuildClassIndex_LookupMissing(t *testing.T) {
	resources := []Resource{
		idxClassResource("r1", "standard"),
	}

	idx := BuildClassIndex(resources)

	got := idx.Lookup("enterprise")
	if len(got) != 0 {
		t.Errorf("expected 0 resources for missing class, got %d", len(got))
	}
}

func TestBuildClassIndex_Classes(t *testing.T) {
	resources := []Resource{
		idxClassResource("r1", "standard"),
		idxClassResource("r2", "premium"),
		idxClassResource("r3", "standard"),
	}

	idx := BuildClassIndex(resources)

	classes := idx.Classes()
	if len(classes) != 2 {
		t.Fatalf("expected 2 distinct classes, got %d", len(classes))
	}

	classSet := make(map[string]bool)
	for _, c := range classes {
		classSet[c] = true
	}
	if !classSet["standard"] {
		t.Errorf("expected 'standard' in classes")
	}
	if !classSet["premium"] {
		t.Errorf("expected 'premium' in classes")
	}
}

func TestBuildClassIndex_EmptyInput(t *testing.T) {
	idx := BuildClassIndex([]Resource{})

	got := idx.Lookup("standard")
	if len(got) != 0 {
		t.Errorf("expected 0 resources from empty index, got %d", len(got))
	}

	classes := idx.Classes()
	if len(classes) != 0 {
		t.Errorf("expected 0 classes from empty index, got %d", len(classes))
	}
}
