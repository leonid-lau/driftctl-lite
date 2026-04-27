package tfstate

import "testing"

// TestProjectFilterAndIndex_RoundTrip verifies that FilterByProject and
// BuildProjectIndex agree on which resources belong to a given project.
func TestProjectFilterAndIndex_RoundTrip(t *testing.T) {
	resources := []Resource{
		projectResource("r1", "platform"),
		projectResource("r2", "platform"),
		projectResource("r3", "data"),
		projectResource("r4", ""),
	}

	const target = "platform"

	filtered := FilterByProject(resources, target, DefaultProjectFilterOptions())
	idx := BuildProjectIndex(resources)
	indexed := idx.Lookup(target)

	if len(filtered) != len(indexed) {
		t.Fatalf("filter returned %d resources but index returned %d",
			len(filtered), len(indexed))
	}

	ids := func(rs []Resource) map[string]bool {
		m := make(map[string]bool, len(rs))
		for _, r := range rs {
			m[r.ID] = true
		}
		return m
	}

	filteredIDs := ids(filtered)
	for _, r := range indexed {
		if !filteredIDs[r.ID] {
			t.Errorf("index contains %q but filter does not", r.ID)
		}
	}
}
