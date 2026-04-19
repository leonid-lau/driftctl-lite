package tfstate

import (
	"testing"
)

func idxZoneResource(id, zone string) Resource {
	return Resource{
		ID:         id,
		Type:       "aws_instance",
		Attributes: map[string]string{"zone": zone},
	}
}

func idxAZResource(id, az string) Resource {
	return Resource{
		ID:         id,
		Type:       "aws_subnet",
		Attributes: map[string]string{"availability_zone": az},
	}
}

func TestBuildZoneIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxZoneResource("r1", "us-east-1a"),
		idxZoneResource("r2", "us-east-1b"),
		idxZoneResource("r3", "us-east-1a"),
	}
	idx := BuildZoneIndex(resources)
	got := idx.Lookup("us-east-1a")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildZoneIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxZoneResource("r1", "US-EAST-1A")}
	idx := BuildZoneIndex(resources)
	got := idx.Lookup("us-east-1a")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildZoneIndex_FallbackToAvailabilityZone(t *testing.T) {
	resources := []Resource{idxAZResource("r1", "eu-west-1b")}
	idx := BuildZoneIndex(resources)
	got := idx.Lookup("eu-west-1b")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildZoneIndex_LookupMissing(t *testing.T) {
	idx := BuildZoneIndex([]Resource{idxZoneResource("r1", "us-west-2a")})
	if got := idx.Lookup("ap-south-1a"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuildZoneIndex_Zones(t *testing.T) {
	resources := []Resource{
		idxZoneResource("r1", "us-east-1a"),
		idxZoneResource("r2", "us-west-2b"),
	}
	idx := BuildZoneIndex(resources)
	if len(idx.Zones()) != 2 {
		t.Fatalf("expected 2 zones, got %d", len(idx.Zones()))
	}
}

func TestBuildZoneIndex_EmptyInput(t *testing.T) {
	idx := BuildZoneIndex(nil)
	if len(idx.Zones()) != 0 {
		t.Fatal("expected no zones")
	}
}
