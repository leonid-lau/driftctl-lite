package tfstate

import (
	"testing"
)

func sevResource(id, severity string) Resource {
	return Resource{
		Type: "aws_instance",
		ID:   id,
		Attributes: map[string]interface{}{"severity": severity},
	}
}

func TestFilterBySeverity_Match(t *testing.T) {
	resources := []Resource{
		sevResource("r1", "high"),
		sevResource("r2", "low"),
		sevResource("r3", "high"),
	}
	opts := DefaultSeverityFilterOptions()
	got := FilterBySeverity(resources, "high", opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterBySeverity_EmptySeverity_ReturnsAll(t *testing.T) {
	resources := []Resource{sevResource("r1", "low"), sevResource("r2", "high")}
	got := FilterBySeverity(resources, "", DefaultSeverityFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterBySeverity_CaseInsensitive(t *testing.T) {
	resources := []Resource{sevResource("r1", "HIGH")}
	got := FilterBySeverity(resources, "high", DefaultSeverityFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterBySeverity_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{sevResource("r1", "HIGH")}
	opts := SeverityFilterOptions{CaseInsensitive: false}
	got := FilterBySeverity(resources, "high", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterBySeverities_ORSemantics(t *testing.T) {
	resources := []Resource{
		sevResource("r1", "low"),
		sevResource("r2", "medium"),
		sevResource("r3", "critical"),
	}
	got := FilterBySeverities(resources, []string{"low", "critical"}, DefaultSeverityFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterBySeverities_EmptyList_ReturnsAll(t *testing.T) {
	resources := []Resource{sevResource("r1", "low"), sevResource("r2", "high")}
	got := FilterBySeverities(resources, nil, DefaultSeverityFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
