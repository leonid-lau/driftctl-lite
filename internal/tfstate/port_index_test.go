package tfstate

import (
	"sort"
	"testing"
)

func idxPortResource(port string) Resource {
	return Resource{
		Type:       "aws_lb_listener",
		ID:         "res-" + port,
		Attributes: map[string]interface{}{"port": port},
	}
}

func TestBuildPortIndex_Lookup(t *testing.T) {
	resources := []Resource{idxPortResource("443"), idxPortResource("80")}
	idx := BuildPortIndex(resources)
	got := idx.Lookup("443")
	if len(got) != 1 || got[0].ID != "res-443" {
		t.Fatalf("expected res-443, got %v", got)
	}
}

func TestBuildPortIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxPortResource("HTTPS")}
	idx := BuildPortIndex(resources)
	got := idx.Lookup("https")
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestBuildPortIndex_LookupMissing(t *testing.T) {
	idx := BuildPortIndex([]Resource{idxPortResource("443")})
	got := idx.Lookup("9999")
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestBuildPortIndex_Ports(t *testing.T) {
	resources := []Resource{idxPortResource("443"), idxPortResource("80"), idxPortResource("22")}
	idx := BuildPortIndex(resources)
	ports := idx.Ports()
	sort.Strings(ports)
	if len(ports) != 3 {
		t.Fatalf("expected 3 ports, got %d: %v", len(ports), ports)
	}
}

func TestBuildPortIndex_EmptyInput(t *testing.T) {
	idx := BuildPortIndex(nil)
	if len(idx.Ports()) != 0 {
		t.Fatal("expected no ports for empty input")
	}
}
