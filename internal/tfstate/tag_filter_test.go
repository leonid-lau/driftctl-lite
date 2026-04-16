package tfstate

import (
	"testing"
)

func tagResource(id, rtype string, attrs map[string]string) Resource {
	return Resource{ID: id, Type: rtype, Attributes: attrs}
}

func TestFilterByTag_MatchKeyAndValue(t *testing.T) {
	resources := []Resource{
		tagResource("r1", "aws_instance", map[string]string{"tag.env": "prod"}),
		tagResource("r2", "aws_instance", map[string]string{"tag.env": "dev"}),
		tagResource("r3", "aws_s3_bucket", map[string]string{"tag.env": "prod"}),
	}
	got := FilterByTag(resources, TagFilter{Key: "env", Value: "prod"})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByTag_MatchKeyOnly(t *testing.T) {
	resources := []Resource{
		tagResource("r1", "aws_instance", map[string]string{"tag.env": "prod"}),
		tagResource("r2", "aws_instance", map[string]string{"tag.team": "ops"}),
	}
	got := FilterByTag(resources, TagFilter{Key: "env"})
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestFilterByTag_EmptyKey_ReturnsAll(t *testing.T) {
	resources := []Resource{
		tagResource("r1", "aws_instance", map[string]string{"tag.env": "prod"}),
		tagResource("r2", "aws_instance", map[string]string{}),
	}
	got := FilterByTag(resources, TagFilter{})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByTags_ANDSemantics(t *testing.T) {
	resources := []Resource{
		tagResource("r1", "aws_instance", map[string]string{"tag.env": "prod", "tag.team": "ops"}),
		tagResource("r2", "aws_instance", map[string]string{"tag.env": "prod", "tag.team": "dev"}),
		tagResource("r3", "aws_instance", map[string]string{"tag.env": "staging"}),
	}
	filters := []TagFilter{
		{Key: "env", Value: "prod"},
		{Key: "team", Value: "ops"},
	}
	got := FilterByTags(resources, filters)
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected only r1, got %v", got)
	}
}

func TestFilterByTag_NoMatch(t *testing.T) {
	resources := []Resource{
		tagResource("r1", "aws_instance", map[string]string{"tag.env": "prod"}),
	}
	got := FilterByTag(resources, TagFilter{Key: "owner", Value: "alice"})
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}
