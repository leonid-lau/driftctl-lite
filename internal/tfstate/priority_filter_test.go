package tfstate

import (
	"testing"
)

func prioResource(id, priority string) Resource {
	return Resource{
		Type: "aws_instance",
		Name: id,
		Attributes: map[string]interface{}{"priority": priority},
	}
}

func TestFilterByPriority_Match(t *testing.T) {
	resources := []Resource{
		prioResource("a", "high"),
		prioResource("b", "low"),
		prioResource("c", "high"),
	}
	got := FilterByPriority(resources, "high", DefaultPriorityFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByPriority_EmptyPriority_ReturnsAll(t *testing.T) {
	resources := []Resource{prioResource("a", "high"), prioResource("b", "low")}
	got := FilterByPriority(resources, "", DefaultPriorityFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByPriority_CaseInsensitive(t *testing.T) {
	resources := []Resource{prioResource("a", "HIGH"), prioResource("b", "low")}
	got := FilterByPriority(resources, "high", DefaultPriorityFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByPriority_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{prioResource("a", "HIGH")}
	opts := PriorityFilterOptions{CaseSensitive: true}
	got := FilterByPriority(resources, "high", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByMinPriority_ReturnsAtOrAbove(t *testing.T) {
	resources := []Resource{
		prioResource("a", "low"),
		prioResource("b", "medium"),
		prioResource("c", "high"),
		prioResource("d", "critical"),
	}
	got := FilterByMinPriority(resources, "high")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByMinPriority_UnknownLevel_ReturnsNil(t *testing.T) {
	resources := []Resource{prioResource("a", "high")}
	got := FilterByMinPriority(resources, "urgent")
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}
