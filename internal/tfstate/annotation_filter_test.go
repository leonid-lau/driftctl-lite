package tfstate

import (
	"testing"
)

func annoResource(id string, annotations map[string]interface{}) Resource {
	attrs := map[string]interface{}{"annotations": annotations}
	return Resource{ID: id, Type: "aws_instance", Attributes: attrs}
}

func TestFilterByAnnotation_MatchKeyAndValue(t *testing.T) {
	resources := []Resource{
		annoResource("r1", map[string]interface{}{"env": "prod"}),
		annoResource("r2", map[string]interface{}{"env": "dev"}),
	}
	got := FilterByAnnotation(resources, "env", "prod", DefaultAnnotationFilterOptions())
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %v", got)
	}
}

func TestFilterByAnnotation_MatchKeyOnly(t *testing.T) {
	resources := []Resource{
		annoResource("r1", map[string]interface{}{"team": "alpha"}),
		annoResource("r2", map[string]interface{}{"other": "x"}),
	}
	got := FilterByAnnotation(resources, "team", "", DefaultAnnotationFilterOptions())
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %v", got)
	}
}

func TestFilterByAnnotation_EmptyKey_ReturnsAll(t *testing.T) {
	resources := []Resource{
		annoResource("r1", map[string]interface{}{"a": "b"}),
		annoResource("r2", map[string]interface{}{"c": "d"}),
	}
	got := FilterByAnnotation(resources, "", "", DefaultAnnotationFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByAnnotations_ANDSemantics(t *testing.T) {
	resources := []Resource{
		annoResource("r1", map[string]interface{}{"env": "prod", "team": "alpha"}),
		annoResource("r2", map[string]interface{}{"env": "prod", "team": "beta"}),
		annoResource("r3", map[string]interface{}{"env": "dev", "team": "alpha"}),
	}
	filters := map[string]string{"env": "prod", "team": "alpha"}
	got := FilterByAnnotations(resources, filters, DefaultAnnotationFilterOptions())
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %v", got)
	}
}

func TestFilterByAnnotation_NoAnnotationsField(t *testing.T) {
	resources := []Resource{
		{ID: "r1", Type: "aws_instance", Attributes: map[string]interface{}{}},
	}
	got := FilterByAnnotation(resources, "env", "prod", DefaultAnnotationFilterOptions())
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}
