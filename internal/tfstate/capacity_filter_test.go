package tfstate

import (
	"testing"
)

func capacityResource(id, capacity string) Resource {
	return Resource{
		Type: "aws_instance",
		ID:   id,
		Attributes: map[string]interface{}{
			"capacity": capacity,
		},
	}
}

func TestFilterByCapacity_Match(t *testing.T) {
	resources := []Resource{
		capacityResource("r1", "large"),
		capacityResource("r2", "small"),
		capacityResource("r3", "large"),
	}
	got := FilterByCapacity(resources, "large", DefaultCapacityFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCapacity_EmptyCapacity_ReturnsAll(t *testing.T) {
	resources := []Resource{
		capacityResource("r1", "large"),
		capacityResource("r2", "small"),
	}
	got := FilterByCapacity(resources, "", DefaultCapacityFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCapacity_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		capacityResource("r1", "Large"),
		capacityResource("r2", "small"),
	}
	got := FilterByCapacity(resources, "large", DefaultCapacityFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByCapacity_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		capacityResource("r1", "Large"),
	}
	opts := CapacityFilterOptions{CaseSensitive: true}
	got := FilterByCapacity(resources, "large", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByCapacities_ORSemantics(t *testing.T) {
	resources := []Resource{
		capacityResource("r1", "large"),
		capacityResource("r2", "small"),
		capacityResource("r3", "medium"),
	}
	got := FilterByCapacities(resources, []string{"large", "medium"}, DefaultCapacityFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCapacities_EmptySlice_ReturnsAll(t *testing.T) {
	resources := []Resource{
		capacityResource("r1", "large"),
		capacityResource("r2", "small"),
	}
	got := FilterByCapacities(resources, nil, DefaultCapacityFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
