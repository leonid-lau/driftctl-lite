package tfstate

import "testing"

func portResource(port string) Resource {
	attrs := map[string]interface{}{}
	if port != "" {
		attrs["port"] = port
	}
	return Resource{Type: "aws_lb_listener", ID: port, Attributes: attrs}
}

func TestFilterByPort_Match(t *testing.T) {
	resources := []Resource{portResource("443"), portResource("80"), portResource("8080")}
	got := FilterByPort(resources, "443", DefaultPortFilterOptions())
	if len(got) != 1 || got[0].ID != "443" {
		t.Fatalf("expected 1 resource with port 443, got %v", got)
	}
}

func TestFilterByPort_EmptyPort_ReturnsAll(t *testing.T) {
	resources := []Resource{portResource("443"), portResource("80")}
	got := FilterByPort(resources, "", DefaultPortFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByPort_CaseInsensitive(t *testing.T) {
	resources := []Resource{portResource("HTTPS"), portResource("http")}
	got := FilterByPort(resources, "https", DefaultPortFilterOptions())
	if len(got) != 1 || got[0].ID != "HTTPS" {
		t.Fatalf("expected 1 match, got %v", got)
	}
}

func TestFilterByPort_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{portResource("HTTPS")}
	opts := PortFilterOptions{CaseSensitive: true}
	got := FilterByPort(resources, "https", opts)
	if len(got) != 0 {
		t.Fatalf("expected no match, got %v", got)
	}
}

func TestFilterByPorts_ORSemantics(t *testing.T) {
	resources := []Resource{portResource("443"), portResource("80"), portResource("22")}
	got := FilterByPorts(resources, []string{"443", "22"}, DefaultPortFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByPorts_Empty_ReturnsAll(t *testing.T) {
	resources := []Resource{portResource("443"), portResource("80")}
	got := FilterByPorts(resources, nil, DefaultPortFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}
