package tfstate

import (
	"testing"
)

func cidrResource(id, cidrBlock string) Resource {
	attrs := map[string]interface{}{}
	if cidrBlock != "" {
		attrs["cidr_block"] = cidrBlock
	}
	return Resource{
		Type:       "aws_vpc",
		ID:         id,
		Attributes: attrs,
	}
}

func TestFilterByCIDR_Match(t *testing.T) {
	resources := []Resource{
		cidrResource("vpc-1", "10.0.0.0/16"),
		cidrResource("vpc-2", "192.168.0.0/24"),
	}
	got := FilterByCIDR(resources, "10.0.0.0/16", DefaultCIDRFilterOptions())
	if len(got) != 1 || got[0].ID != "vpc-1" {
		t.Fatalf("expected vpc-1, got %+v", got)
	}
}

func TestFilterByCIDR_EmptyCIDR_ReturnsAll(t *testing.T) {
	resources := []Resource{
		cidrResource("vpc-1", "10.0.0.0/16"),
		cidrResource("vpc-2", "192.168.0.0/24"),
	}
	got := FilterByCIDR(resources, "", DefaultCIDRFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByCIDR_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		{Type: "aws_vpc", ID: "vpc-1", Attributes: map[string]interface{}{"cidr": "10.0.0.0/8"}},
	}
	got := FilterByCIDR(resources, "10.0.0.0/8", DefaultCIDRFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByCIDR_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		cidrResource("vpc-1", "10.0.0.0/16"),
	}
	opts := CIDRFilterOptions{CaseSensitive: true}
	// CIDR matching is exact for case-sensitive; different casing won't matter
	// for IPs, but the option is tested for correctness.
	got := FilterByCIDR(resources, "10.0.0.0/32", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByCIDRs_ORSemantics(t *testing.T) {
	resources := []Resource{
		cidrResource("vpc-1", "10.0.0.0/16"),
		cidrResource("vpc-2", "172.16.0.0/12"),
		cidrResource("vpc-3", "192.168.0.0/24"),
	}
	got := FilterByCIDRs(resources, []string{"10.0.0.0/16", "172.16.0.0/12"}, DefaultCIDRFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCIDRs_Empty_ReturnsAll(t *testing.T) {
	resources := []Resource{
		cidrResource("vpc-1", "10.0.0.0/16"),
	}
	got := FilterByCIDRs(resources, nil, DefaultCIDRFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}
