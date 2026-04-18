package tfstate

import (
	"testing"
)

func zoneResource(id, zone string) Resource {
	return Resource{
		ID:   id,
		Type: "aws_instance",
		Attributes: map[string]interface{}{
			"zone": zone,
		},
	}
}

func TestFilterByZone_Match(t *testing.T) {
	resources := []Resource{
		zoneResource("r1", "us-east-1a"),
		zoneResource("r2", "us-west-2b"),
	}
	got := FilterByZone(resources, "us-east-1a", DefaultZoneFilterOptions())
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %v", got)
	}
}

func TestFilterByZone_EmptyZone_ReturnsAll(t *testing.T) {
	resources := []Resource{zoneResource("r1", "us-east-1a"), zoneResource("r2", "eu-west-1b")}
	got := FilterByZone(resources, "", DefaultZoneFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByZone_CaseInsensitive(t *testing.T) {
	resources := []Resource{zoneResource("r1", "US-EAST-1A")}
	got := FilterByZone(resources, "us-east-1a", DefaultZoneFilterOptions())
	if len(got) != 1 {
		t.Fatal("expected case-insensitive match")
	}
}

func TestFilterByZone_CaseSensitive_NoMatch(t *testing.T) {
	opts := ZoneFilterOptions{CaseSensitive: true}
	resources := []Resource{zoneResource("r1", "US-EAST-1A")}
	got := FilterByZone(resources, "us-east-1a", opts)
	if len(got) != 0 {
		t.Fatal("expected no match with case-sensitive option")
	}
}

func TestFilterByZones_ORSemantics(t *testing.T) {
	resources := []Resource{
		zoneResource("r1", "us-east-1a"),
		zoneResource("r2", "eu-west-1b"),
		zoneResource("r3", "ap-southeast-1a"),
	}
	got := FilterByZones(resources, []string{"us-east-1a", "eu-west-1b"}, DefaultZoneFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByZones_EmptyList_ReturnsAll(t *testing.T) {
	resources := []Resource{zoneResource("r1", "us-east-1a")}
	got := FilterByZones(resources, []string{}, DefaultZoneFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}
