package tfstate

import (
	"testing"
)

func tagResource(id, tagKey, tagVal string) Resource {
	attrs := map[string]string{}
	if tagKey != "" {
		attrs["tags."+tagKey] = tagVal
	}
	return Resource{ID: id, Type: "aws_instance", Attributes: attrs}
}

func TestFilterByTag_MatchKeyAndValue(t *testing.T) {
	resources := []Resource{
		tagResource("r1", "env", "prod"),
		tagResource("r2", "env", "dev"),
		tagResource("r3", "team", "platform"),
	}
	got := FilterByTag(resources, "env", "prod")
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %+v", got)
	}
}

func TestFilterByTag_MatchKeyOnly(t *testing.T) {
	resources := []Resource{
		tagResource("r1", "env", "prod"),
		tagResource("r2", "env", "dev"),
		tagResource("r3", "team", "platform"),
	}
	got := FilterByTag(resources, "env", "")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByTag_EmptyKey_ReturnsAll(t *testing.T) {
	resources := []Resource{
		tagResource("r1", "env", "prod"),
		tagResource("r2", "team", "sre"),
	}
	got := FilterByTag(resources, "", "")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByTags_ANDSemantics(t *testing.T) {
	r1 := Resource{ID: "r1", Type: "aws_instance", Attributes: map[string]string{
		"tags.env": "prod", "tags.team": "sre",
	}}
	r2 := Resource{ID: "r2", Type: "aws_instance", Attributes: map[string]string{
		"tags.env": "prod",
	}}
	resources := []Resource{r1, r2}
	got := FilterByTags(resources, map[string]string{"env": "prod", "team": "sre"})
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %+v", got)
	}
}

func TestFilterByTags_EmptyMap_ReturnsAll(t *testing.T) {
	resources := []Resource{
		tagResource("r1", "env", "prod"),
		tagResource("r2", "env", "dev"),
	}
	got := FilterByTags(resources, map[string]string{})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
