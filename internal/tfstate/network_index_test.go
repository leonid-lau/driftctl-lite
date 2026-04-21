package tfstate

import (
	"sort"
	"testing"
)

func idxNetworkResource(network string) Resource {
	return Resource{
		Type:       "aws_subnet",
		Name:       "r",
		Attributes: map[string]interface{}{"network": network},
	}
}

func TestBuildNetworkIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxNetworkResource("prod-net"),
		idxNetworkResource("dev-net"),
	}
	idx := BuildNetworkIndex(resources)
	got := idx.Lookup("prod-net")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildNetworkIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxNetworkResource("PROD-NET")}
	idx := BuildNetworkIndex(resources)
	got := idx.Lookup("prod-net")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildNetworkIndex_LookupMissing(t *testing.T) {
	idx := BuildNetworkIndex([]Resource{idxNetworkResource("prod-net")})
	got := idx.Lookup("missing")
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestBuildNetworkIndex_Networks(t *testing.T) {
	resources := []Resource{
		idxNetworkResource("prod-net"),
		idxNetworkResource("dev-net"),
		idxNetworkResource("prod-net"),
	}
	idx := BuildNetworkIndex(resources)
	networks := idx.Networks()
	sort.Strings(networks)
	if len(networks) != 2 || networks[0] != "dev-net" || networks[1] != "prod-net" {
		t.Fatalf("unexpected networks: %v", networks)
	}
}

func TestBuildNetworkIndex_FallbackToVPC(t *testing.T) {
	r := Resource{
		Type:       "aws_instance",
		Name:       "r",
		Attributes: map[string]interface{}{"vpc": "vpc-abc"},
	}
	idx := BuildNetworkIndex([]Resource{r})
	got := idx.Lookup("vpc-abc")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildNetworkIndex_EmptyInput(t *testing.T) {
	idx := BuildNetworkIndex(nil)
	if len(idx.Networks()) != 0 {
		t.Fatal("expected no networks")
	}
}
