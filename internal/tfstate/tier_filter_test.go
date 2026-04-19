package tfstate

import (
	"testing"
)

func tierResource(id, tier string) Resource {
	attrs := map[string]interface{}{}
	if tier != "" {
		attrs["tier"] = tier
	}
	return Resource{ID: id, Type: "aws_instance", Attributes: attrs}
}

func TestFilterByTier_Match(t *testing.T) {
	resources := []Resource{
		tierResource("a", "frontend"),
		tierResource("b", "backend"),
		tierResource("c", "frontend"),
	}
	got := FilterByTier(resources, "frontend", DefaultTierFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByTier_EmptyTier_ReturnsAll(t *testing.T) {
	resources := []Resource{tierResource("a", "frontend"), tierResource("b", "backend")}
	got := FilterByTier(resources, "", DefaultTierFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByTier_CaseInsensitive(t *testing.T) {
	resources := []Resource{tierResource("a", "Frontend")}
	got := FilterByTier(resources, "frontend", DefaultTierFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByTier_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{tierResource("a", "Frontend")}
	opts := TierFilterOptions{CaseSensitive: true}
	got := FilterByTier(resources, "frontend", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByTiers_ORSemantics(t *testing.T) {
	resources := []Resource{
		tierResource("a", "frontend"),
		tierResource("b", "backend"),
		tierResource("c", "data"),
	}
	got := FilterByTiers(resources, []string{"frontend", "data"}, DefaultTierFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildTierIndex_Lookup(t *testing.T) {
	resources := []Resource{tierResource("a", "frontend"), tierResource("b", "backend")}
	idx := BuildTierIndex(resources)
	got := idx.Lookup("frontend")
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestBuildTierIndex_LookupMissing(t *testing.T) {
	idx := BuildTierIndex([]Resource{tierResource("a", "frontend")})
	if got := idx.Lookup("data"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuildTierIndex_Tiers(t *testing.T) {
	resources := []Resource{tierResource("a", "frontend"), tierResource("b", "backend")}
	idx := BuildTierIndex(resources)
	if len(idx.Tiers()) != 2 {
		t.Fatalf("expected 2 tiers, got %d", len(idx.Tiers()))
	}
}
