package tfstate

import (
	"testing"
)

func costResource(id, costCenter string) Resource {
	attrs := map[string]interface{}{}
	if costCenter != "" {
		attrs["cost_center"] = costCenter
	}
	return Resource{ID: id, Type: "aws_instance", Attributes: attrs}
}

func TestFilterByCost_Match(t *testing.T) {
	resources := []Resource{
		costResource("r1", "engineering"),
		costResource("r2", "marketing"),
		costResource("r3", ""),
	}
	got := FilterByCost(resources, "engineering", DefaultCostFilterOptions())
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %v", got)
	}
}

func TestFilterByCost_EmptyCost_ReturnsAll(t *testing.T) {
	resources := []Resource{costResource("r1", "eng"), costResource("r2", "ops")}
	got := FilterByCost(resources, "", DefaultCostFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCost_CaseInsensitive(t *testing.T) {
	resources := []Resource{costResource("r1", "Engineering")}
	got := FilterByCost(resources, "engineering", DefaultCostFilterOptions())
	if len(got) != 1 {
		t.Fatal("expected match")
	}
}

func TestFilterByCost_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{costResource("r1", "Engineering")}
	opts := CostFilterOptions{CaseSensitive: true}
	got := FilterByCost(resources, "engineering", opts)
	if len(got) != 0 {
		t.Fatal("expected no match")
	}
}

func TestFilterByCosts_ORSemantics(t *testing.T) {
	resources := []Resource{
		costResource("r1", "engineering"),
		costResource("r2", "marketing"),
		costResource("r3", "finance"),
	}
	got := FilterByCosts(resources, []string{"engineering", "finance"}, DefaultCostFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildCostIndex_Lookup(t *testing.T) {
	resources := []Resource{
		costResource("r1", "engineering"),
		costResource("r2", "marketing"),
	}
	idx := BuildCostIndex(resources)
	result := idx.Lookup("Engineering")
	if len(result) != 1 || result[0].ID != "r1" {
		t.Fatalf("expected r1, got %v", result)
	}
}

func TestBuildCostIndex_LookupMissing(t *testing.T) {
	idx := BuildCostIndex([]Resource{costResource("r1", "eng")})
	if idx.Lookup("finance") != nil {
		t.Fatal("expected nil")
	}
}
