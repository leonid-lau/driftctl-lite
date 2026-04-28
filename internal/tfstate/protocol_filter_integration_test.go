package tfstate

import "testing"

// TestProtocolFilterAndIndex_RoundTrip verifies that BuildProtocolIndex and
// FilterByProtocol agree on membership for a mixed set of resources.
func TestProtocolFilterAndIndex_RoundTrip(t *testing.T) {
	resources := []Resource{
		protocolResource("tcp"),
		protocolResource("udp"),
		protocolResource("icmp"),
		protocolResource("TCP"), // duplicate with different case
	}

	const target = "tcp"

	filtered := FilterByProtocol(resources, target, DefaultProtocolFilterOptions())
	idx := BuildProtocolIndex(resources)
	indexed := idx.Lookup(target)

	if len(filtered) != len(indexed) {
		t.Errorf("filter returned %d resources, index returned %d",
			len(filtered), len(indexed))
	}

	ids := make(map[string]bool)
	for _, r := range filtered {
		ids[r.ID] = true
	}
	for _, r := range indexed {
		if !ids[r.ID] {
			t.Errorf("index returned resource %q not found in filter results", r.ID)
		}
	}
}
