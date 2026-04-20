package tfstate

import (
	"testing"
)

// TestDependencyFilterAndIndex_RoundTrip verifies that BuildDependencyIndex
// and FilterByDependency agree on which resources belong to a given dependency.
func TestDependencyFilterAndIndex_RoundTrip(t *testing.T) {
	resources := []Resource{
		idxDepResource("core"),
		idxDepResource("core"),
		idxDepResource("networking"),
		idxDepResource("storage"),
	}

	target := "core"
	opts := DefaultDependencyFilterOptions()

	filtered := FilterByDependency(resources, target, opts)
	idx := BuildDependencyIndex(resources)
	looked := idx.Lookup(target)

	if len(filtered) != len(looked) {
		t.Fatalf("filter returned %d, index returned %d", len(filtered), len(looked))
	}

	for i := range filtered {
		fAttr := filtered[i].Attributes["dependency"]
		lAttr := looked[i].Attributes["dependency"]
		if fAttr != lAttr {
			t.Errorf("mismatch at index %d: filter=%v index=%v", i, fAttr, lAttr)
		}
	}
}
