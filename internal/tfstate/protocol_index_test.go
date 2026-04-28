package tfstate

import (
	"testing"
)

func idxProtocolResource(protocol string) Resource {
	return Resource{
		Type: "aws_security_group_rule",
		ID:   "sgr-" + protocol,
		Attributes: map[string]string{
			"protocol": protocol,
		},
	}
}

func TestBuildProtocolIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxProtocolResource("tcp"),
		idxProtocolResource("udp"),
	}
	idx := BuildProtocolIndex(resources)

	got := idx.Lookup("tcp")
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].ID != "sgr-tcp" {
		t.Errorf("unexpected resource ID: %s", got[0].ID)
	}
}

func TestBuildProtocolIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxProtocolResource("TCP")}
	idx := BuildProtocolIndex(resources)

	if got := idx.Lookup("tcp"); len(got) != 1 {
		t.Errorf("expected 1 result for case-insensitive lookup, got %d", len(got))
	}
	if got := idx.Lookup("TCP"); len(got) != 1 {
		t.Errorf("expected 1 result for upper-case lookup, got %d", len(got))
	}
}

func TestBuildProtocolIndex_LookupMissing(t *testing.T) {
	idx := BuildProtocolIndex([]Resource{idxProtocolResource("tcp")})
	if got := idx.Lookup("icmp"); got != nil {
		t.Errorf("expected nil for missing protocol, got %v", got)
	}
}

func TestBuildProtocolIndex_Protocols(t *testing.T) {
	resources := []Resource{
		idxProtocolResource("tcp"),
		idxProtocolResource("udp"),
		idxProtocolResource("tcp"), // duplicate
	}
	idx := BuildProtocolIndex(resources)
	protocols := idx.Protocols()
	if len(protocols) != 2 {
		t.Errorf("expected 2 distinct protocols, got %d", len(protocols))
	}
}

func TestBuildProtocolIndex_EmptyInput(t *testing.T) {
	idx := BuildProtocolIndex(nil)
	if got := idx.Lookup("tcp"); got != nil {
		t.Errorf("expected nil on empty index, got %v", got)
	}
	if p := idx.Protocols(); len(p) != 0 {
		t.Errorf("expected empty protocols slice, got %v", p)
	}
}
