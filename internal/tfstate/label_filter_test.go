package tfstate

import (
	"testing"
)

func labelResource(id, rtype string, labels map[string]string) Resource {
	return Resource{ID: id, Type: rtype, Labels: labels}
}

func TestFilterByLabel_MatchKeyAndValue(t *testing.T) {
	resources := []Resource{
		labelResource("a", "aws_s3_bucket", map[string]string{"env": "prod"}),
		labelResource("b", "aws_s3_bucket", map[string]string{"env": "dev"}),
	}
	got := FilterByLabel(resources, "env", "prod")
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("expected 1 result with id=a, got %+v", got)
	}
}

func TestFilterByLabel_MatchKeyOnly(t *testing.T) {
	resources := []Resource{
		labelResource("a", "aws_s3_bucket", map[string]string{"env": "prod"}),
		labelResource("b", "aws_instance", map[string]string{"region": "us-east-1"}),
	}
	got := FilterByLabel(resources, "env", "")
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("expected 1 result, got %+v", got)
	}
}

func TestFilterByLabel_EmptyKey_ReturnsAll(t *testing.T) {
	resources := []Resource{
		labelResource("a", "aws_s3_bucket", map[string]string{"env": "prod"}),
		labelResource("b", "aws_instance", nil),
	}
	got := FilterByLabel(resources, "", "")
	if len(got) != 2 {
		t.Fatalf("expected all resources, got %d", len(got))
	}
}

func TestFilterByLabels_ANDSemantics(t *testing.T) {
	resources := []Resource{
		labelResource("a", "aws_s3_bucket", map[string]string{"env": "prod", "team": "infra"}),
		labelResource("b", "aws_s3_bucket", map[string]string{"env": "prod"}),
	}
	opts := DefaultLabelFilterOptions()
	got := FilterByLabels(resources, map[string]string{"env": "prod", "team": "infra"}, opts)
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("expected 1 AND match, got %+v", got)
	}
}

func TestFilterByLabels_ORSemantics(t *testing.T) {
	resources := []Resource{
		labelResource("a", "aws_s3_bucket", map[string]string{"env": "prod"}),
		labelResource("b", "aws_instance", map[string]string{"team": "infra"}),
		labelResource("c", "aws_vpc", map[string]string{"region": "us-west-2"}),
	}
	opts := LabelFilterOptions{RequireAll: false}
	got := FilterByLabels(resources, map[string]string{"env": "prod", "team": "infra"}, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2 OR matches, got %d", len(got))
	}
}
