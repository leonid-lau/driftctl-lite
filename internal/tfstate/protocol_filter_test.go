package tfstate

import (
	"testing"
)

func protocolResource(id, protocol string) Resource {
	return Resource{
		ID:   id,
		Type: "aws_lb_listener",
		Attributes: map[string]interface{}{
			"protocol": protocol,
		},
	}
}

func TestFilterByProtocol_Match(t *testing.T) {
	resources := []Resource{
		protocolResource("r1", "HTTPS"),
		protocolResource("r2", "HTTP"),
		protocolResource("r3", "TCP"),
	}
	got := FilterByProtocol(resources, "HTTPS", DefaultProtocolFilterOptions())
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected [r1], got %v", got)
	}
}

func TestFilterByProtocol_EmptyProtocol_ReturnsAll(t *testing.T) {
	resources := []Resource{
		protocolResource("r1", "HTTPS"),
		protocolResource("r2", "HTTP"),
	}
	got := FilterByProtocol(resources, "", DefaultProtocolFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByProtocol_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		protocolResource("r1", "https"),
		protocolResource("r2", "HTTP"),
	}
	got := FilterByProtocol(resources, "HTTPS", DefaultProtocolFilterOptions())
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected [r1], got %v", got)
	}
}

func TestFilterByProtocol_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		protocolResource("r1", "https"),
	}
	opts := FilterOptions{CaseSensitive: true}
	got := FilterByProtocol(resources, "HTTPS", opts)
	if len(got) != 0 {
		t.Fatalf("expected no results, got %v", got)
	}
}

func TestFilterByProtocols_ORSemantics(t *testing.T) {
	resources := []Resource{
		protocolResource("r1", "HTTPS"),
		protocolResource("r2", "HTTP"),
		protocolResource("r3", "TCP"),
	}
	got := FilterByProtocols(resources, []string{"HTTPS", "TCP"}, DefaultProtocolFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByProtocols_EmptySlice_ReturnsAll(t *testing.T) {
	resources := []Resource{
		protocolResource("r1", "HTTPS"),
		protocolResource("r2", "HTTP"),
	}
	got := FilterByProtocols(resources, []string{}, DefaultProtocolFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}
