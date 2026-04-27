package tfstate

import (
	"testing"
)

func poolResource(pool string) Resource {
	attrs := map[string]string{}
	if pool != "" {
		attrs["pool"] = pool
	}
	return Resource{Type: "aws_instance", ID: pool + "-id", Attributes: attrs}
}

func TestFilterByPool_Match(t *testing.T) {
	resources := []Resource{
		poolResource("alpha"),
		poolResource("beta"),
		poolResource("alpha"),
	}
	got := FilterByPool(resources, "alpha", DefaultPoolFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByPool_EmptyPool_ReturnsAll(t *testing.T) {
	resources := []Resource{poolResource("alpha"), poolResource("beta")}
	got := FilterByPool(resources, "", DefaultPoolFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByPool_CaseInsensitive(t *testing.T) {
	resources := []Resource{poolResource("Alpha"), poolResource("BETA")}
	got := FilterByPool(resources, "alpha", DefaultPoolFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByPool_CaseSensitive_NoMatch(t *testing.T) {
	opts := PoolFilterOptions{CaseSensitive: true}
	resources := []Resource{poolResource("Alpha")}
	got := FilterByPool(resources, "alpha", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByPools_ORSemantics(t *testing.T) {
	resources := []Resource{
		poolResource("alpha"),
		poolResource("beta"),
		poolResource("gamma"),
	}
	got := FilterByPools(resources, []string{"alpha", "gamma"}, DefaultPoolFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildPoolIndex_Lookup(t *testing.T) {
	resources := []Resource{poolResource("alpha"), poolResource("beta"), poolResource("alpha")}
	idx := BuildPoolIndex(resources)
	got := idx.Lookup("alpha")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildPoolIndex_LookupMissing(t *testing.T) {
	idx := BuildPoolIndex([]Resource{poolResource("alpha")})
	if idx.Lookup("missing") != nil {
		t.Fatal("expected nil for missing pool")
	}
}

func TestBuildPoolIndex_Pools(t *testing.T) {
	resources := []Resource{poolResource("alpha"), poolResource("beta")}
	idx := BuildPoolIndex(resources)
	pools := idx.Pools()
	if len(pools) != 2 {
		t.Fatalf("expected 2 pools, got %d", len(pools))
	}
}
