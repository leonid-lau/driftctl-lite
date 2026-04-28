package tfstate

import (
	"testing"

	"github.com/snyk/driftctl-lite/internal/tfstate/resource"
)

func endpointResource(id, endpoint string) resource.Resource {
	return resource.Resource{
		ID:   id,
		Type: "aws_lb",
		Attributes: map[string]string{
			"endpoint": endpoint,
		},
	}
}

func TestFilterByEndpoint_Match(t *testing.T) {
	resources := []resource.Resource{
		endpointResource("r1", "https://api.example.com"),
		endpointResource("r2", "https://other.example.com"),
	}
	got := FilterByEndpoint(resources, "https://api.example.com", DefaultEndpointFilterOptions())
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %v", got)
	}
}

func TestFilterByEndpoint_EmptyEndpoint_ReturnsAll(t *testing.T) {
	resources := []resource.Resource{
		endpointResource("r1", "https://api.example.com"),
		endpointResource("r2", "https://other.example.com"),
	}
	got := FilterByEndpoint(resources, "", DefaultEndpointFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByEndpoint_CaseInsensitive(t *testing.T) {
	resources := []resource.Resource{
		endpointResource("r1", "HTTPS://API.EXAMPLE.COM"),
	}
	got := FilterByEndpoint(resources, "https://api.example.com", DefaultEndpointFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByEndpoint_CaseSensitive_NoMatch(t *testing.T) {
	resources := []resource.Resource{
		endpointResource("r1", "HTTPS://API.EXAMPLE.COM"),
	}
	opts := EndpointFilterOptions{CaseSensitive: true}
	got := FilterByEndpoint(resources, "https://api.example.com", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByEndpoints_ORSemantics(t *testing.T) {
	resources := []resource.Resource{
		endpointResource("r1", "https://api.example.com"),
		endpointResource("r2", "https://other.example.com"),
		endpointResource("r3", "https://third.example.com"),
	}
	got := FilterByEndpoints(resources,
		[]string{"https://api.example.com", "https://other.example.com"},
		DefaultEndpointFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByEndpoint_MissingAttribute(t *testing.T) {
	resources := []resource.Resource{
		{ID: "r1", Type: "aws_lb", Attributes: map[string]string{}},
	}
	got := FilterByEndpoint(resources, "https://api.example.com", DefaultEndpointFilterOptions())
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}
